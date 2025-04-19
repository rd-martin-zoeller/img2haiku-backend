package compose

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func makePrompt(language string, tags []string) (string, *types.ComposeError) {
	var prompt string

	data := struct {
		Language   string
		TagsString string
	}{
		Language:   language,
		TagsString: makeTagsString(tags),
	}

	const promptTemplate = `
	First: Check if the image is appropriate. If the image violates any policy or includes disallowed content (e.g., explicit violence, sexual content, hate symbols, etc.), do not generate a Haiku. Instead, return a JSON error message such as:
        {
          "error": "<Error message in {{.Language}}>"
        }

        If the image is appropriate, proceed as follows:

        1. Consider the user-provided tags that describe the image’s emotion, mood, or atmosphere: {{.TagsString}}. Make use of these tags to guide your interpretation and the overall tone of the resulting Haiku.
        2. Look at the image from a human’s perspective and describe it in one concise sentence, identifying the general subject and emotional tone. Keep the aforementioned tags in mind. If no tags are provided, use your own best judgment.
        3. Using both the user-supplied tags and the concise description, generate a Haiku in {{.Language}} with the following rules:
            - Exactly three lines
            - No rhyming
            - Short and evocative
            - Incorporate sensory imagery and the emotion/mood from the tags
        4. If the image’s content is unclear, focus on a single visible element (e.g., color, light, or shapes) and the feeling it evokes.
        5. Return the final answer in valid JSON with the following structure:
        {
          "description": "<one-sentence description of the image in {{.Language}}>",
          "haiku": "<the three-line poem in {{.Language}}>"
        }
        6. Do not wrap the final JSON answer in markdown or any other formatting.
        7. Do not include any explanations, disclaimers, or additional keys beyond "description" and "haiku" in the JSON output.
	`

	template, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return prompt, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to parse prompt template: " + err.Error(),
		}
	}

	var buff bytes.Buffer
	if err := template.Execute(&buff, data); err != nil {
		return prompt, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to execute prompt template: " + err.Error(),
		}
	}

	return buff.String(), nil
}

func makeTagsString(tags []string) string {
	if len(tags) == 0 {
		return "No tags provided"
	}

	return strings.Join(tags, ", ")
}
