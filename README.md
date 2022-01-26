# 開発
12/24~開発中  
メルカリのようなサイト
top  
![](/src/app/static/img/ui/top.png)
product-page  
![](/src/app/static/img/ui/product.png)
マイページいろいろ  
![](/src/app/static/img/ui/hist.png)



# 機能
ユーザー認証  
    - ログイン  
        - recaptcha v3  
    - ユーザー登録  
        - メール認証  

商品登録  
    - stripeAPIで商品登録  
    - 登録商品一覧表示  

決済  
    - stripeAPIで連結アカウント作成  
    - プラットフォームに振込  
    - クライアントに振込  

その他  
    - 購入者情報表示機能  
    - 返品先情報表示機能  
    - 発送、到着確認  
    - etc  
  


# 環境
コンテナ  
    - docker  
バックエンド  
    - golang  
    - postgresql
フロントエンド  
    - html  
    - css  
    - js  
AWS  
    - s3  
        - 未設定  
    - ecs  
        - 未設定  
    - ecr  
        - 未設定  
    - rds  
        - 未設定  
    - route53  
        - 未設定  
CI/CD  
    - CircleCI  
        - 未設定  

決済システム  
    - stripe    
  
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
publish_key = xxx  
eps = xxx  

[recaptcha]  
pk = xxx   