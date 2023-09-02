resource "edison_talk" "sdk_team" {
  title            = "A Modern Terraform Plugin Framework"
  description      = "An introduction to a new framework for developing Terraform plugins."
  duration_minutes = 30
  prerecorded      = true
  speaker_ids = [
    edison_speaker.katy.id,
    edison_speaker.brian.id,
    edison_speaker.paddy.id,
  ]
}

output "brian_recording" {
  value = edison_talk.sdk_team.recordings["Brian Flad"]
}
