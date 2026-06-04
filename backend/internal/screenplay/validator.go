package screenplay

import "fmt"

func Validate(input Screenplay) ValidationResult {
	var errors []string

	if input.Title == "" {
		errors = append(errors, "title is required")
	}
	if len(input.Characters) == 0 {
		errors = append(errors, "characters are required")
	}
	if len(input.Scenes) == 0 {
		errors = append(errors, "scenes are required")
	}

	for i, scene := range input.Scenes {
		label := fmt.Sprintf("scenes[%d]", i)
		if scene.ID == "" {
			errors = append(errors, label+".id is required")
		}
		if scene.Location == "" {
			errors = append(errors, label+".location is required")
		}
		if scene.Time == "" {
			errors = append(errors, label+".time is required")
		}
		if scene.Summary == "" {
			errors = append(errors, label+".summary is required")
		}

		for j, dialogue := range scene.Dialogues {
			dialogueLabel := fmt.Sprintf("%s.dialogues[%d]", label, j)
			if dialogue.Character == "" {
				errors = append(errors, dialogueLabel+".character is required")
			}
			if dialogue.Line == "" {
				errors = append(errors, dialogueLabel+".line is required")
			}
		}
	}

	return ValidationResult{
		Passed: len(errors) == 0,
		Errors: errors,
	}
}
