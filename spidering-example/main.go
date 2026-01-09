package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"regexp"

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
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UA)
	if ref != "" {
		req.Header.Set("Referer", ref)
	}
	res, _ := cl.Do(req)
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	return doc
}

// 2. あなたが選んだ「短い getSession」
func getSession() (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// 認証ページ取得
	doc := doReq(client, AuthURL, "")

	// リンク抽出
	link, _ := doc.Find("a.btn-ageauth-yes").Attr("href")
	if link == "" {
		return nil, fmt.Errorf("認証ボタンが見つかりませんでした")
	}

	// 認証実行（Referer付き）
	doReq(client, link, AuthURL)

	return client, nil
}

// 3. 実行部分
func main() {
	client, err := getSession()
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	// あとはこの client を使って好きなページを叩くだけ
	fmt.Println("✅ 認証完了。データを取得します...")

	doc := doReq(client, "https://www.sokmil.com/idol/_item/item511685.htm", "")
	fmt.Println("タイトル:", doc.Find("title").Text())
	fmt.Println("H1の内容:", doc.Find("h1").Text())

	// 1. HTML全体を文字列として取得
	html, _ := doc.Html()

	// 2. 正規表現で video_url: '...' の中身を探す
	// パターン：video_url: '(ここを抜き出す)'
	re := regexp.MustCompile(`video_url:\s*'(https?://[^']+)'`)
	match := re.FindStringSubmatch(html)

	if len(match) > 1 {
		videoURL := match[1]
		fmt.Println("動画URL取得成功:", videoURL)
	} else {
		fmt.Println("動画URLが見つかりませんでした")
	}
}
