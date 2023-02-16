# TikTok Microservices

<p align="center"> 
  <img src="https://upload.wikimedia.org/wikipedia/en/thumb/a/a9/TikTok_logo.svg/1200px-TikTok_logo.svg.png" alt="tiktok logo"/> 
</p>

This repository contains the coursework for the module ECM3408 (Enterprise Computing). It contains three microservices: 

 * `tracks` - A microservice that controls access to a tracks database, consisting of an ID (track name) and Audio (base64 encoded `.wav` file);

 * `search` - A service that contains a client for [AudD](https://audd.io), a music recognition software similar to [Shazam](https://www.shazam.com/gb/home);

 * `cooltown` - A "facade" service that uses the other two microservices to request a base64 encoded audio file of a hummed tune, and responds with a base64 encoded audio file of the song in question.

## Deployment

From the repository directory, you can run the following:

```bash 
(cd addison/tracks && go run main.go) & \
(cd addison/search && go run main.go) & \
(cd addison/cooltown && go run main.go) &
```

## Testing

To test these service, the database MUST be deleted prior to running. After
deleting the database and starting the three services, you can run the tests
with:

```bash 
python3 testing/e2e_test.py
```
