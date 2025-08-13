миграции:
goose create init sql    из директории где создать файл миграции
goose postgres "host=127.0.0.1 port=5435 user=test password=test dbname=users sslmode=disable" up        (down / status)
goose create insert_test_polls sql  
после отмены неудачной миграции удаляем файл с миграцией