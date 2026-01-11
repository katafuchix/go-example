package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	// go get github.com/go-resty/resty/v"
)

// ブラウザのふりをするための共通User-Agent
const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

/*
func getSession() (*http.Client, error) {
	// 1. クッキーを自動管理する設定
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// 2. 年齢認証ページへアクセス（ブラウザを装う）
	req, _ := http.NewRequest("GET", "https://www.sokmil.com/member/ageauth/", nil)
	req.Header.Set("User-Agent", ua)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 3. HTML解析して「はい」ボタンのURLを抽出
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// 確実に取るために「btn-ageauth-yes」クラスを含むaタグを検索
	ageOkURL, exists := doc.Find("a.btn-ageauth-yes").Attr("href")
	if !exists {
		// 取れなかった場合はHTMLをデバッグ出力して終了
		html, _ := doc.Html()
		fmt.Println("--- DEBUG HTML ---")
		fmt.Println(html)
		return nil, fmt.Errorf("認証ボタンが見つかりませんでした")
	}

	// 4. 認証確定リンクへアクセス（PythonのReferer設定を再現）
	reqOk, _ := http.NewRequest("GET", ageOkURL, nil)
	reqOk.Header.Set("User-Agent", ua)
	reqOk.Header.Set("Referer", "https://www.sokmil.com/member/ageauth/")

	resOk, err := client.Do(reqOk)
	if err != nil {
		return nil, err
	}
	defer resOk.Body.Close()

	fmt.Println("✅ 認証成功: ", resOk.Request.URL.String())
	return client, nil
}
*/

// 共通設定
const (
	UA      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	AuthURL = "https://www.sokmil.com/member/ageauth/"
)

// 1. あなたが選んだ「短い doReq」
func doReq(cl *http.Client, url, ref string) *goquery.Document {
	// 新しいHTTP GETリクエストを作成
	req, _ := http.NewRequest("GET", url, nil)
	// ヘッダーにユーザーエージェントを設定
	req.Header.Set("User-Agent", UA)
	// もし参照元のURLが指定されている場合は、Refererヘッダーを設定
	if ref != "" {
		req.Header.Set("Referer", ref)
	}
	// リクエストを実行し、レスポンスを取得
	res, _ := cl.Do(req)
	// 関数終了時にレスポンスのボディを閉じるようにdeferを使う
	defer res.Body.Close()
	// レスポンスボディからgoquery用のドキュメントを作成
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	// ドキュメントを返す
	return doc
}

// 2. あなたが選んだ「短い getSession」
func getSession() (*http.Client, error) {
	// クッキージャーを新規に作成
	jar, _ := cookiejar.New(nil)
	// HTTPクライアントをクッキージャーと共に作成
	client := &http.Client{Jar: jar}

	// 認証ページ取得
	doc := doReq(client, AuthURL, "")

	// リンク抽出
	link, _ := doc.Find("a.btn-ageauth-yes").Attr("href")
	if link == "" {
		// リンクが見つからなかった場合はエラーを返す
		return nil, fmt.Errorf("認証ボタンが見つかりませんでした")
	}

	// 認証実行（Referer付き）
	// 抽出したリンクに対して、Refererを付けて認証リクエストを送信
	doReq(client, link, AuthURL)

	// 認証完了したクライアントを返す
	return client, nil
}

// 3. 実行部分
func main() {
	client, err := getSession()
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	fmt.Println("✅ 認証完了。データを取得します...")

	// 1. 一覧ページからアイテムのURLをごっそり取得
	fmt.Println("一覧ページを解析中...")
	pageURLs := getItemURLs(client, "https://www.sokmil.com/idol/")
	fmt.Printf("%d 個のアイテムが見つかりました\n", len(pageURLs))

	// 2. 並列処理の準備
	var wg sync.WaitGroup
	limit := make(chan struct{}, 3) // 同時3つまでに制限

	for i, url := range pageURLs {
		wg.Add(1)
		go func(target string, index int) {
			defer wg.Done()

			limit <- struct{}{} // 枠に入る
			fmt.Printf("[%d / %d] 処理中...\n", index+1, len(pageURLs))
			crawlAndDownload(client, target) // ページ解析〜保存まで一気に実行
			<-limit                          // 枠を空ける
		}(url, i)
	}

	wg.Wait()
	fmt.Println("すべての処理が完了しました！")
}

func getItemURLs(cl *http.Client, listURL string) []string {
	// HTTPリクエストを送信してレスポンスのHTMLドキュメントを取得
	doc := doReq(cl, listURL, "")
	// アイテムのURLを格納するスライスを宣言
	var urls []string

	// aタグのhref属性をすべてチェック
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		// 条件： "/idol/_item/item" を含んでいるものだけを抽出
		if exists && strings.Contains(href, "/idol/_item/item") {
			// もし相対パス（/idol/..）なら絶対パスに変換
			if strings.HasPrefix(href, "/") {
				href = "https://www.sokmil.com" + href
			}
			// URLをスライスに追加
			urls = append(urls, href)
		}
	})

	// 重複を排除する（同じリンクが複数ある場合があるため）
	return unique(urls)
}

// 重複を消すための便利関数
func unique(slice []string) []string {
	m := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !m[s] {
			m[s] = true
			result = append(result, s)
		}
	}
	return result
}

func crawlAndDownload(cl *http.Client, pageURL string) {
	// 1. ページの中身を取得
	doc := doReq(cl, pageURL, "") // 以前作った共通関数を使用
	html, _ := doc.Html()

	// 2. 正規表現で動画URL (video_url) を抜き出す
	re := regexp.MustCompile(`video_url:\s*'(https?://[^']+)'`)
	match := re.FindStringSubmatch(html)

	if len(match) < 2 {
		fmt.Printf("動画URLが見つかりませんでした: %s\n", pageURL)
		return
	}
	videoURL := match[1]

	// --- 保存先の設定 ---
	saveDir := "mp4"

	// 1. mp4フォルダがなければ作成する (0755は一般的な権限設定)
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.Mkdir(saveDir, 0755)
	}

	// 3. 動画URLからファイル名を抜き出す (f=xxxx.mp4)
	reFile := regexp.MustCompile(`f=([^&]+)`)
	fileMatch := reFile.FindStringSubmatch(videoURL)

	fileName := ""
	if len(fileMatch) > 1 {
		fileName = fileMatch[1]
	} else {
		// ID部分を抽出（例: item511685）
		reID := regexp.MustCompile(`item\d+`)
		idMatch := reID.FindString(pageURL)
		if idMatch != "" {
			fileName = idMatch + ".mp4"
		} else {
			// 最悪のケース：ランダムな値を付ける（衝突防止）
			fileName = fmt.Sprintf("unknown_%d.mp4", time.Now().UnixNano())
		}
	}

	// 「mp4/ファイル名」というパスを作る
	// filepath.Joinを使うと、WindowsでもMacでも正しくスラッシュを処理してくれます
	savePath := filepath.Join(saveDir, fileName)

	// 4. ダウンロード実行
	out, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer out.Close()

	// HTTP GETリクエストを作成
	req, _ := http.NewRequest("GET", videoURL, nil)

	// リクエストヘッダーにUser-Agentを設定
	req.Header.Set("User-Agent", UA)

	// HTTPクライアント(cl)でリクエストを送信し、レスポンスを取得
	resp, err := cl.Do(req)
	if err != nil {
		// エラーが発生した場合は処理を中断
		return
	}
	// プログラムの終了時にレスポンスボディを必ず閉じるようにす
	defer resp.Body.Close()

	fmt.Printf("開始: %s\n", fileName)
	io.Copy(out, resp.Body)
	fmt.Printf("完了: %s\n", fileName)
}
