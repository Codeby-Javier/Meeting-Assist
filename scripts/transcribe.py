import speech_recognition as sr
import sys
import os

def transcribe(audio_path, language="id-ID"):
    r = sr.Recognizer()
    
    # Check if file exists
    if not os.path.exists(audio_path):
        print(f"Error: File not found: {audio_path}")
        sys.exit(1)

    try:
        with sr.AudioFile(audio_path) as source:
            audio_data = r.record(source)
            # Use Google Speech Recognition (free tier)
            text = r.recognize_google(audio_data, language=language)
            print(text)
    except sr.UnknownValueError:
        print("[Unintelligible]")
    except sr.RequestError as e:
        print(f"Error: Could not request results; {e}")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python transcribe.py <audio_path> [language_code]")
        sys.exit(1)
    
    path = sys.argv[1]
    lang = sys.argv[2] if len(sys.argv) > 2 else "id-ID"
    transcribe(path, lang)
