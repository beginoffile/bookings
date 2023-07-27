go build -o bookings.exe ./cmd/web
bookings.exe -dbname=bookings -dbuser=postgres -dbpass="12345678" -cache=false -production=false