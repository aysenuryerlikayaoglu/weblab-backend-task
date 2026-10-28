# Etkinlik Kayıt & Kontenjan Yönetimi API

Bu proje, bellek içi (in-memory) çalışan bir RESTful etkinlik yönetim servisidir. Go ve Gin framework kullanılarak geliştirilmiş olup, harici bir veritabanı bağlantısı gerektirmez.

## Çalıştırma Talimatları

Projeyi yerel ortamınızda çalıştırmak için aşağıdaki komutu kullanabilirsiniz:

go run main.go

Bu komut ile sunucu varsayılan olarak `http://localhost:8080` portunda ayağa kalkacaktır.

## Endpointler

* **GET /events**: Kayıtlı etkinlikleri listeler. `?city=Istanbul` şeklinde parametre alabilir.
* **POST /events**: Yeni etkinlik ekler. JSON body: `{"id":"1", "title":"Konser", "city":"Istanbul", "capacity":50}`
* **POST /events/{id}/register**: Belirtilen etkinliğe kayıt yapar ve kontenjanı 1 azaltır.