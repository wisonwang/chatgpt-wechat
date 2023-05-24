package stablediffusion

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func Text2Image(host, key, prompt string) (string, error) {
	url := host + "/api/v3/text2img"
	reqPayload := map[string]interface{}{
		"prompt":              prompt,
		"key":                 key,
		"negative_prompt":     "((out of frame)), ((extra fingers)), mutated hands, ((poorly drawn hands)), ((poorly drawn face)), (((mutation))), (((deformed))), (((tiling))), ((naked)), ((tile)), ((fleshpile)), ((ugly)), (((abstract))), blurry, ((bad anatomy)), ((bad proportions)), ((extra limbs)), cloned face, (((skinny))), glitchy, ((extra breasts)), ((double torso)), ((extra arms)), ((extra hands)), ((mangled fingers)), ((missing breasts)), (missing lips), ((ugly face)), ((fat)), ((extra legs)), anime",
		"width":               "512",
		"height":              "512",
		"samples":             "1",
		"num_inference_steps": "20",
		"safety_checker":      "no",
		"enhance_prompt":      "yes",

		"guidance_scale": 7.5,
	}

	client := &http.Client{}
	body, _ := json.Marshal(reqPayload)
	drawReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		logx.Info("draw request client build fail", err)
		return "", err
	}
	logx.Info("draw request client build success")
	drawReq.Header.Add("Content-Type", "application/json")
	res, err := client.Do(drawReq)
	if err != nil {
		logx.Info("draw request fail", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logx.Info("draw request fail", err)

		return "", err
	}

	var resPayload map[string]interface{}
	err = json.Unmarshal(resBody, &resPayload)
	if err != nil {
		logx.Info("draw request fail", err)
		return "", err
	}
	if resPayload["status"] != nil && resPayload["status"].(string) == "success" && resPayload["output"] != nil {
		images := resPayload["output"].([]interface{})
		for _, image := range images {
			s := image.(string)
			return s, err
		}
	}

	return "", errors.New("生成图片失败")
}
