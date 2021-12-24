# tables

- user
    - id　(PRIMARY KEY)
    - username
    - password　(hash)
    - email

- accouts
    - id　(PRIMARY KEY)
    - user_id　(userテーブルのid)
    - stripe_account

- products
    - id　(PRIMARY KEY)
    - user_id　(販売者、userテーブルのid)
    - product_name
    - amount
    - quantity　(数量)

- settlement
    - id　(PRIMARY KEY)
    - user_id　(購入者、userテーブルのid)
    - product_id　(productsテーブルのid)
