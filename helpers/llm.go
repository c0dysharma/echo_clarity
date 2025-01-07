package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/c0dysharma/echo_clarity/structs"
	"github.com/charmbracelet/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

const modelName = "gemini-1.5-flash"

func CallLLM(imagePath string)([]structs.CalendarEvent, error){
	llmPrompt := fmt.Sprintf(`You are a programmer and productive person you write down the tasks you need to do in your day along with the time range you are going to do start time and end time in UTC+05:30 timezone of today's date that is %s and the task you need to, given an image you need to OCR the data and return the data in following JSON format not anything else

The format
[
	{
	startTime : <time in RFC3339 format>
	endTime: <time in RFC3339 format>
	eventName: <string>
	}
]`, time.Now().Format("02/01/2006")) 
	ctx := context.Background()
  apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

  llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
  if err != nil {
    log.Fatal(err)
  }

	imgData, err := os.ReadFile(imagePath)
  if err != nil {
    log.Fatal(err)
  }

  parts := []llms.ContentPart{
    llms.BinaryPart("image/png", imgData),
    llms.TextPart(llmPrompt),
  }

  content := []llms.MessageContent{
    {
      Role:  llms.ChatMessageTypeHuman,
      Parts: parts,
    },
  }

  resp, err := llm.GenerateContent(ctx, content, llms.WithModel(modelName))
  if err != nil {
    log.Fatal(err)
  }

  if len(resp.Choices) > 0 {
      content := resp.Choices[0].Content
			cleanedContent, err := ExtractJSONString(content)

			if err != nil {
				log.Error("Error extracting JSON string", "Error", err)
				return []structs.CalendarEvent{}, err
			}

			// unmarshal the content in events
			var events []structs.CalendarEvent
			err = json.Unmarshal([]byte(cleanedContent), &events)

			if err!= nil {
        log.Error("Error unmarshalling events", "Error", err)
        return []structs.CalendarEvent{}, err
      }

      return events, nil
  } else {
      return []structs.CalendarEvent{}, nil
  }

}