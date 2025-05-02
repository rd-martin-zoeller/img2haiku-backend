package compose

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/utils"
)

const promptTemplateWithoutTags = `
First: Check if the image is appropriate. If it violates any policy, ignore the rest of this prompt, and instead, return:
{
	"error": "<A descriptive error message in {{.Language}}>"
}

Otherwise, proceed as follows:

1. The output of this task should really impress the user by how well it captures the essence of the image. Here are some guidelines to help you:
	- If the image has a funny or silly subject, be funny and silly.
	- If the image has a serious or dramatic subject, use a significantly more serious tone.
	- If you feel like the image is a work of art, be poetic and artistic.
	- If you feel like the image captures an important memory, be heartfelt and emotional.
	- And so on.
2. Look at the image from a human’s perspective. Infer what makes this image interesting to the user by asking yourself the following questions: What is the subject of the image? What is happening in the image? What is the mood or emotion conveyed by the image? What are the colors, shapes, and textures present in the image?
3. Using the answers to the above questions, describe the image in one or two sentences. Use the following language: {{.Language}}. Keep it concise, without leaving out any important details. Remember your description.
4. If the image’s content is unclear, focus on a single visible element (e.g., color, light, or shapes) and the feeling it evokes.
5. Using everything you have learned about the image so far, generate a Haiku in {{.Language}} with the following rules and guidelines:
	- Exactly three lines, separated by \\n
	- No rhyming
	- Do not use dashes or colons
	- Make it powerful and evocative
	- Don't be abstract or vague
	- Don't be afraid to use strong imagery or metaphors to convey the essence of the image, without exaggerating or making it absurd
5. Return the final answer in valid JSON with the following structure:
{
	"description": "<one-sentence description of the image in {{.Language}}>",
	"haiku": "<the three-line poem in {{.Language}}>"
}
6. Do not wrap the final JSON answer in markdown or any other formatting.
7. Do not include any explanations, disclaimers, or additional keys beyond "description" and "haiku" in the JSON output.
`

const promptTemplateWithTags = `
First: Check if the image is appropriate. If it violates any policy, ignore the rest of this prompt, and instead, return:
{
	"error": "<A descriptive error message in {{.Language}}>"
}

Otherwise, proceed as follows:
The user has provided the following tags that they say are relevant to the image: {{.TagsString}}. Use these tags to help you understand the image better and to generate the output.
1. The output of this task should really impress the user by how well it captures the essence of the image in accordance to the tags they provided.
2. Look at the image from a human’s perspective. Infer what makes this image interesting to the user by asking yourself the following questions: What is the subject of the image? What is happening in the image? How do the provided tags fit the image? What are the colors, shapes, and textures present in the image?
3. Using the answers to the above questions, describe the image in one or two sentences. Use the following language: {{.Language}}. Keep it concise, without leaving out any important details. Remember your description.
4. If the image’s content is unclear, focus on a single visible element (e.g., color, light, or shapes) and the feeling it evokes.
5. Using everything you have learned about the image so far, generate a Haiku in {{.Language}} with the following rules and guidelines:
	- Exactly three lines, separated by \\n
	- No rhyming
	- Do not use dashes or colons
	- Make it powerful and evocative
	- Don't be abstract or vague
	- Don't be afraid to use strong imagery or metaphors to convey the essence of the image, without exaggerating or making it absurd
5. Return the final answer in valid JSON with the following structure:
{
	"description": "<one-sentence description of the image in {{.Language}}>",
	"haiku": "<the three-line poem in {{.Language}}>"
}
6. Do not wrap the final JSON answer in markdown or any other formatting.
7. Do not include any explanations, disclaimers, or additional keys beyond "description" and "haiku" in the JSON output.
`

func makePrompt(language string, tags []string) (string, error) {
	var prompt string

	data := struct {
		Language   string
		TagsString string
	}{
		Language:   language,
		TagsString: makeTagsString(tags),
	}

	promptTemplate := pickTemplate(tags)

	template, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return prompt, utils.NewInternalErr("Failed to parse prompt template: %s", err.Error())
	}

	var buff bytes.Buffer
	if err := template.Execute(&buff, data); err != nil {
		return prompt, utils.NewInternalErr("Failed to execute prompt template: %s", err.Error())
	}

	return buff.String(), nil
}

func makeTagsString(tags []string) string {
	if len(tags) == 0 {
		return "No tags provided"
	}

	return strings.Join(tags, ", ")
}

func pickTemplate(tags []string) string {
	if len(tags) == 0 {
		return promptTemplateWithoutTags
	}

	return promptTemplateWithTags
}
