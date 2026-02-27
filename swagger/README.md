из корня запускаем:

swag init -g ./core/cmd/main.go -o ./swagger/core
swag init --dir ./msvc_auth -o ./swagger/auth

специально ограничиваемся через --dir, а не используем -g иначе почему-то перемешиваются комментарии у сервисов
ТАК НЕ ДЕЛАЕМ: swag init -g ./msvc_auth/main.go -o ./swagger/auth