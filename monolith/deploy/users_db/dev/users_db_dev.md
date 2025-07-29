docker-compose up -d
если меняем данные в compose, то не забываем удалять volume postgres, иначе подрузится volume с прежней базой


для production:
сертификаты, pg_hba и пр можно без Dockerfile через compose сделать.
файлы для деплоя на проде, включая секреты, держим отдельно от гита