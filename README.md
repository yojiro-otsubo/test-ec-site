# 開発
12/24~開発中  
メルカリのようなサイト

![](/src/app/static/img/ui/top.png)

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
コンテナ  
    - docker  
バックエンド  
    - golang  
フロントエンド  
    - html  
    - css  
    - js  
AWS  
    - s3  
    - ecs  
    - ecr  
    - rds  
    - route53  
CI/CD  
    - CircleCI  
  
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