# 開発
12/24~開発中  
メルカリのようなサイト

# 機能
ユーザー認証  
    - ログイン  
    - ユーザー登録  
  
商品登録  
    - stripeAPIで商品登録  
    - 登録商品一覧表示  

決済  
    - stripeAPIで連結アカウント作成   
    - プラットフォームに振込  
    - クライアントに振込
    - 振込タイミング  
  


# 環境
Golang  
js  
  
# 決済システム
stripe  
  
# config.ini
  
[web]  
port = xxxx  
  
[db]  
driver = xxx  
db_host = xxx  
name = xxx  
user = xxx  
password = xxx  
  
[stripe]  
stripe_key = xxx  