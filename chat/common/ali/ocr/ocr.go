package ocr

import (
	"net/http"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"

	// 1、这里只是以ocr为例，其他能力请引入相应类目的包。包名可参考本文档上方的SDK包名称，能力名可参考对应API文档中的Action参数。例如您想使用通用分割，其文档为https://help.aliyun.com/document_detail/151960.html，可以知道该能力属于分割抠图类目，能力名称为SegmentCommonImage，那么您需要将代码中ocr20191230改为imageseg20191230，ocr-20191230改为imageseg-20191230，将RecognizeBankCard改为SegmentCommonImage。版本可以参考文档上面的SDK包信息。
	ocr20191230 "github.com/alibabacloud-go/ocr-20191230/v3/client"

	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *ocr20191230.Client, _err error) {
	config := &openapi.Config{
		// 2、"YOUR_ACCESS_KEY_ID", "YOUR_ACCESS_KEY_SECRET" 的生成请参考https://help.aliyun.com/document_detail/175144.html
		// 如果您是用的子账号AccessKey，还需要为子账号授予权限AliyunVIAPIFullAccess，请参考https://help.aliyun.com/document_detail/145025.html
		// 您的 AccessKey ID
		AccessKeyId: tea.String(*accessKeyId),
		// 您的 AccessKey Secret
		AccessKeySecret: tea.String(*accessKeySecret),
	}

	// 3、访问的域名。注意：这个地方需要求改为相应类目的域名，参考：https://help.aliyun.com/document_detail/143103.html
	config.Endpoint = tea.String("ocr.cn-shanghai.aliyuncs.com")
	// 4、这里只是以ocr为例，其他能力请使用相应类目的包
	client, err := ocr20191230.NewClient(config)
	return client, err
}

func Image2Txt(arg string, client *ocr20191230.Client) (res string, _err error) {
	httpClient := http.Client{}
	file, _ := httpClient.Get(arg)

	recognizeGeneralRequest := &ocr20191230.RecognizeCharacterAdvanceRequest{
		ImageURLObject:    file.Body,
		MinHeight:         tea.Int32(10),
		OutputProbability: tea.Bool(false),
	}
	runtime := &util.RuntimeOptions{}

	txt, tryErr := func() (res string, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		result, _err := client.RecognizeCharacterAdvance(recognizeGeneralRequest, runtime)
		if _err != nil {
			return "", _err
		}

		for _, l := range result.Body.Data.Results {
			res += *(l.Text) + "\n"
		}

		return res, nil
	}()

	if tryErr != nil {
		var sdkError = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			sdkError = _t
		} else {
			sdkError.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		r, _err := util.AssertAsString(sdkError.Message)
		if _err != nil {
			return *r, _err
		}
	}
	return txt, _err
}
