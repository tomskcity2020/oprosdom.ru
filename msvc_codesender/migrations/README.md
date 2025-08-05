миграции:
goose create init sql    из директории где создать файл миграции
goose postgres "host=127.0.0.1 port=5432 user=test password=test dbname=notify sslmode=disable" up        (down / status)
goose create alter_sms_messages sql  
после отмены неудачной миграции удаляем файл с миграцией