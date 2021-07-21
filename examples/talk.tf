resource "hashitalks_talk" "sdk_team" {
  title            = "A Modern Terraform Plugin Framework"
  description      = "An introduction to a new framework for developing Terraform plugins."
  duration_minutes = 30
  prerecorded      = true
  speaker_ids = [
    hashitalks_speaker.katy.id,
    hashitalks_speaker.brian.id,
    hashitalks_speaker.paddy.id,
  ]
}

output "brian_recording" {
  value = hashitalks_talk.sdk_team.recordings["Brian Flad"]
}
