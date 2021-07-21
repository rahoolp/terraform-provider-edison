resource "hashitalks_workshop" "provider" {
  title            = "Building Your First Provider"
  description      = "Learn to build your very first provider with the new Terraform Plugin Framework."
  duration_minutes = 120
  presenters = {
    "Luke" = {
      title    = "Education Engineer"
      employer = "DadCorp"
      pronouns = "he/him"
    },
    "Elliot" = {
      title    = "Lead Instructor"
      employer = "HashiCorp"
      pronouns = "they/them"
    },
    "Serene" = {
      title    = "Terraform Expert"
      pronouns = "she/her"
    },
  }
  meeting_info = {
    url      = "https://zoom.us/my/tfsdkteam"
    password = "super secret zoom password!"
  }
}
