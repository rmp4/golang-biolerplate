# 刪除原有的 git 紀錄
rm -rf .git

# 初始化新的 git 倉庫
git init

# 刪除或清空 go.sum
rm go.sum

# 更新依賴
go mod tidy
