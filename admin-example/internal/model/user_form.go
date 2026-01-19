package model

// UserCreateForm: 新規登録用
type UserCreateForm struct {
	// 修正箇所: max=20" (閉じ) + 半角スペース + label
	Name  string `form:"name" validate:"required,min=3,max=20" label:"名前"`
	Name1 string `form:"name1" validate:"required,min=3,max=20" label:"名前1"`
}

// UserUpdateForm: 更新用
type UserUpdateForm struct {
	ID   uint64 `form:"id" validate:"required" label:"ID"`
	Name string `form:"name" validate:"required,min=3,max=20" label:"名前"`
}
