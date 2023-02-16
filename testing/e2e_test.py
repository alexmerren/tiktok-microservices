import requests
import base64 
from urllib.parse import quote

SEARCH_PORT = 3001
TRACKS_PORT = 3000
COOLTOWN_PORT = 3002

SEARCH_API_URL = f"http://localhost:{SEARCH_PORT}/search"
TRACKS_API_URL = f"http://localhost:{TRACKS_PORT}/tracks"
COOLTOWN_API_URL = f"http://localhost:{COOLTOWN_PORT}/cooltown"

RECORDING_FILENAMES = [
        "Everybody+(Backstreet's+Back)+(Radio+Edit).wav",
        "Blinding+Lights.wav",
        "Don't+Look+Back+In+Anger.wav",
        "good+4+u.wav"
        ]

HUMMING_FILENAMES = [
        "~Everybody+(Backstreet's+Back)+(Radio+Edit).wav",
        "~Blinding+Lights.wav",
        "~Don't+Look+Back+In+Anger.wav",
        "~good+4+u.wav"
        ]

DUMMY_HUMMING_FILENAME = "soundfile.wav"

def encode_file_to_base_64(filename: str) -> str:
    return str(base64.b64encode(open(filename, "rb").read()))

def url_escape_string(string: str) -> str:
    return quote(string)

def search_test() -> None:
    def search_200() -> None:
        for track in HUMMING_FILENAMES:
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            request_data = {
                    "audio": encoded_audio
                    }
            response = requests.post(SEARCH_API_URL, json=request_data)
            assert(response.status_code == 200)

    def search_400() -> None:
        request_data = {}
        response = requests.post(SEARCH_API_URL, json=request_data)
        assert(response.status_code == 400)

    def search_404() -> None:
        encoded_audio = encode_file_to_base_64(DUMMY_HUMMING_FILENAME)[2:-1]
        request_data = {
                "audio": encoded_audio
                }
        response = requests.post(SEARCH_API_URL, json=request_data)
        assert(response.status_code == 404)

    def search_500() -> None:
        pass

    search_200()
    search_400()
    search_404()

def tracks_test() -> None:
    def create_201() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            request_id = track[:-4]
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            request_data = {
                    "id": request_id,
                    "audio": encoded_audio,
                    }
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 201)
            response = requests.delete(track_url)
            assert(response.status_code == 204)

    def create_204() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            request_id = track[:-4]
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            request_data = {
                    "id": request_id,
                    "audio": encoded_audio,
                    }
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 201)
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 204)
            response = requests.delete(track_url)
            assert(response.status_code == 204)

    def create_400() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            request_id = track[:-4]
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            request_data = {
                    "audio": encoded_audio
                    }
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 400)
            request_data = {
                    "id": request_id,
                    }
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 400)
            malformed_track_url = f"{TRACKS_API_URL}/incorrect-id"
            request_data = {
                    "id": request_id,
                    }
            response = requests.put(malformed_track_url, json=request_data)
            assert(response.status_code == 400)

    def list_200() -> None:
        response = requests.get(TRACKS_API_URL)
        assert(response.status_code == 200)

    def read_200() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            request_id = track[:-4]
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            request_data = {
                    "id": request_id,
                    "audio": encoded_audio
                    }
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 201)
            response = requests.get(track_url)
            assert(response.status_code == 200)
            response = requests.delete(track_url)
            assert(response.status_code == 204)

    def read_404() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            response = requests.get(track_url)
            assert(response.status_code == 404)


    def delete_204() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            request_id = track[:-4]
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            request_data = {
                    "id": request_id,
                    "audio": encoded_audio
                    }
            response = requests.put(track_url, json=request_data)
            assert(response.status_code == 201)
            response = requests.delete(track_url)
            assert(response.status_code == 204)

    def delete_404() -> None:
        for track in RECORDING_FILENAMES:
            encoded_id = url_escape_string(track[:-4])
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            response = requests.delete(track_url)
            assert(response.status_code == 404)

    def create_500() -> None:
        pass

    def list_500() -> None:
        pass

    def read_500() -> None:
        pass

    def delete_500() -> None:
        pass

    create_201()
    create_204()
    create_400()
    list_200()
    read_200()
    read_404()
    delete_204()
    delete_404()

def cooltown_test() -> None:
    def cooltown_200() -> None:
        for humming, track in zip(HUMMING_FILENAMES, RECORDING_FILENAMES):
            encoded_id = url_escape_string(track[:-4])
            request_id = track[:-4]
            encoded_audio = encode_file_to_base_64(track)[2:-1]
            track_url = f"{TRACKS_API_URL}/{encoded_id}"
            tracks_request_data = {
                    "id": request_id,
                    "audio": encoded_audio,
                    }

            response = requests.put(track_url, json=tracks_request_data)
            assert(response.status_code == 201)

            encoded_audio = encode_file_to_base_64(humming)[2:-1]
            cooltown_request_data = {
                    "audio": encoded_audio,
                    }
            response = requests.post(COOLTOWN_API_URL, json=cooltown_request_data)
            assert(response.status_code == 200)

            response = requests.delete(track_url)
            assert(response.status_code == 204)

    def cooltown_400() -> None:
        cooltown_request_data = {}
        response = requests.post(COOLTOWN_API_URL, json=cooltown_request_data)
        assert(response.status_code == 400)

        cooltown_request_data = {
                "audio": 12312381628736187236,
                }
        response = requests.post(COOLTOWN_API_URL, json=cooltown_request_data)
        assert(response.status_code == 400)

    def cooltown_404() -> None:
        encoded_audio = encode_file_to_base_64(DUMMY_HUMMING_FILENAME)[2:-1]
        cooltown_request_data = {
                "audio": encoded_audio,
                }
        response = requests.post(COOLTOWN_API_URL, json=cooltown_request_data)
        assert(response.status_code == 404)

    def cooltown_500() -> None:
        pass

    cooltown_200()
    cooltown_400()
    cooltown_404()

def main() -> None:
    tracks_test()
    search_test()
    cooltown_test()

if __name__ == "__main__":
    main()
