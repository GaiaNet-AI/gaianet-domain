// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vhost

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/util/version"
)

var NotFoundPagePath = ""

const (
	NotFound = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0" />
    <title>Gaia</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        background: #f5f5ef;
      }
      .root {
        width: 100%;
        max-width: 1440px;
        margin: 0 auto;
        height: 100vh;
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-family: sans-serif;
        gap: 24px;
      }
      .error-desktop {
        display: flex;
        align-items: center;
        gap: 24px;
      }
      .error-mobile {
        display: none;
        align-items: center;
        gap: 10px;
      }
      .info {
        display: flex;
        flex-direction: column;
        padding-right: 40px;
      }
      .info h2 {
        color: #121314;
        font-size: 32px;
        font-weight: 500;
      }

      .info p {
        color: #666;
        font-size: 17px;
        margin-top: 20px;
      }

      .info a {
        color: #121314;
      }
      @media (max-width: 1024px) {
        .error-desktop {
          display: none;
        }
        .error-mobile {
          display: flex;
        }
      }
      @media (max-width: 768px) {
        .root {
          flex-direction: column;
          align-items: start;
          justify-content: center;
          font-family: sans-serif;
          gap: 84px;
        }
        .info {
          padding-right: 20px;
          padding-left: 20px;
        }

        .info h2 {
          font-size: 20px;
          font-weight: 500;
        }

        .info p {
          color: #666;
          font-size: 14px;
          margin-top: 14px;
        }
      }
    </style>
  </head>
  <body>
    <div class="root">
      <div class="error-desktop">
        <svg xmlns="http://www.w3.org/2000/svg" width="260" height="299" viewBox="0 0 260 299" fill="none">
          <path
            d="M189.524 206.753V298.948H145.729V206.753H-22V184.441L140.138 0.0517578H189.524V171.391H259.877V206.753H189.524ZM145.729 171.391V47.2016L37.6368 171.391H145.729Z"
            fill="#121314"
          />
        </svg>
        <svg xmlns="http://www.w3.org/2000/svg" width="232" height="309" viewBox="0 0 232 309" fill="none">
          <path
            d="M116.077 0C133.471 0 149.312 2.52589 163.6 7.57766C177.888 12.3488 190.001 19.5055 199.941 29.0477C210.191 38.3093 217.956 49.6758 223.236 63.1471C228.827 76.6185 231.623 91.9142 231.623 109.034V199.966C231.623 217.928 228.827 233.785 223.236 247.537C217.646 261.008 209.725 272.375 199.475 281.636C189.225 290.617 176.956 297.493 162.668 302.264C148.69 306.755 133.16 309 116.077 309C98.9931 309 83.3074 306.755 69.0194 302.264C55.042 297.493 42.9283 290.617 32.6782 281.636C22.4281 272.375 14.5076 261.008 8.91665 247.537C3.3257 233.785 0.530212 217.928 0.530212 199.966V109.034C0.530212 91.9142 3.17039 76.6185 8.45074 63.1471C14.0417 49.6758 21.8069 38.3093 31.7464 29.0477C41.9964 19.5055 54.2655 12.3488 68.5535 7.57766C82.8415 2.52589 98.6825 0 116.077 0ZM116.077 271.954C127.88 271.954 138.13 270.129 146.827 266.481C155.834 262.832 163.289 257.781 169.191 251.326C175.092 244.871 179.441 237.293 182.236 228.593C185.342 219.612 186.895 209.929 186.895 199.545V109.455C186.895 99.0708 185.342 89.5286 182.236 80.8283C179.441 71.8474 175.092 64.1294 169.191 57.6744C163.289 51.2193 155.834 46.1676 146.827 42.5191C138.13 38.8706 127.88 37.0463 116.077 37.0463C104.273 37.0463 93.8681 38.8706 84.8604 42.5191C76.1634 46.1676 68.8641 51.2193 62.9625 57.6744C57.061 64.1294 52.5572 71.8474 49.4511 80.8283C46.6556 89.5286 45.2578 99.0708 45.2578 109.455V199.545C45.2578 209.929 46.6556 219.612 49.4511 228.593C52.5572 237.293 57.061 244.871 62.9625 251.326C68.8641 257.781 76.1634 262.832 84.8604 266.481C93.8681 270.129 104.273 271.954 116.077 271.954Z"
            fill="#121314"
          />
        </svg>
        <svg xmlns="http://www.w3.org/2000/svg" width="283" height="299" viewBox="0 0 283 299" fill="none">
          <path
            d="M212.018 206.753V298.948H168.222V206.753H0.493713V184.441L162.631 0.0517578H212.018V171.391H282.371V206.753H212.018ZM168.222 171.391V47.2016L60.1305 171.391H168.222Z"
            fill="#121314"
          />
        </svg>
      </div>

      <div class="error-mobile">
        <svg xmlns="http://www.w3.org/2000/svg" width="86" height="116" viewBox="0 0 86 116" fill="none">
          <path
            d="M58.7763 79.8633V115.07H42.0517V79.8633H-22V71.3429L39.9166 0.928711H58.7763V66.3592H85.6424V79.8633H58.7763ZM42.0517 66.3592V18.9342L0.773932 66.3592H42.0517Z"
            fill="#121314"
          />
        </svg>
        <svg xmlns="http://www.w3.org/2000/svg" width="90" height="118" viewBox="0 0 90 118" fill="none">
          <path
            d="M45.0353 0C51.6777 0 57.727 0.964578 63.1832 2.89373C68.6395 4.71572 73.2655 7.44868 77.0611 11.0926C80.9754 14.6294 83.9407 18.97 85.9572 24.1144C88.0922 29.2589 89.1598 35.0999 89.1598 41.6376V76.3624C89.1598 83.2216 88.0922 89.277 85.9572 94.5286C83.8221 99.673 80.7974 104.014 76.8832 107.55C72.9689 110.98 68.2836 113.606 62.8274 115.428C57.4897 117.143 51.5591 118 45.0353 118C38.5115 118 32.5215 117.143 27.0652 115.428C21.7276 113.606 17.1016 110.98 13.1873 107.55C9.27308 104.014 6.24841 99.673 4.11336 94.5286C1.9783 89.277 0.910767 83.2216 0.910767 76.3624V41.6376C0.910767 35.0999 1.91899 29.2589 3.93543 24.1144C6.07049 18.97 9.03584 14.6294 12.8315 11.0926C16.7458 7.44868 21.431 4.71572 26.8873 2.89373C32.3435 0.964578 38.3929 0 45.0353 0ZM45.0353 103.853C49.5426 103.853 53.4569 103.156 56.7781 101.763C60.2179 100.37 63.0646 98.4405 65.3183 95.9755C67.572 93.5104 69.2326 90.6167 70.3001 87.2943C71.4862 83.8647 72.0793 80.1671 72.0793 76.2016V41.7984C72.0793 37.8329 71.4862 34.1889 70.3001 30.8665C69.2326 27.4369 67.572 24.4896 65.3183 22.0245C63.0646 19.5595 60.2179 17.6303 56.7781 16.2371C53.4569 14.8438 49.5426 14.1471 45.0353 14.1471C40.5279 14.1471 36.5544 14.8438 33.1145 16.2371C29.7933 17.6303 27.0059 19.5595 24.7522 22.0245C22.4986 24.4896 20.7787 27.4369 19.5925 30.8665C18.525 34.1889 17.9912 37.8329 17.9912 41.7984V76.2016C17.9912 80.1671 18.525 83.8647 19.5925 87.2943C20.7787 90.6167 22.4986 93.5104 24.7522 95.9755C27.0059 98.4405 29.7933 100.37 33.1145 101.763C36.5544 103.156 40.5279 103.853 45.0353 103.853Z"
            fill="#121314"
          />
        </svg>
        <svg xmlns="http://www.w3.org/2000/svg" width="109" height="116" viewBox="0 0 109 116" fill="none">
          <path
            d="M81.4332 79.8633V115.07H64.7085V79.8633H0.65686V71.3429L62.5735 0.928711H81.4332V66.3592H108.299V79.8633H81.4332ZM64.7085 66.3592V18.9342L23.4308 66.3592H64.7085Z"
            fill="#121314"
          />
        </svg>
      </div>

      <div class="info">
        <h2>GaiaNet node is not started</h2>
        <p>
          Visit
          <a href="https://www.gaianet.ai" target="_blank">https://www.gaianet.ai</a>
          for more information
        </p>
      </div>
    </div>
  </body>
</html>
`
)

func getNotFoundPageContent() []byte {
	var (
		buf []byte
		err error
	)
	if NotFoundPagePath != "" {
		buf, err = os.ReadFile(NotFoundPagePath)
		if err != nil {
			log.Warnf("read custom 404 page error: %v", err)
			buf = []byte(NotFound)
		}
	} else {
		buf = []byte(NotFound)
	}
	return buf
}

func NotFoundResponse() *http.Response {
	header := make(http.Header)
	header.Set("server", "frp/"+version.Full())
	header.Set("Content-Type", "text/html")

	content := getNotFoundPageContent()
	res := &http.Response{
		Status:        "Not Found",
		StatusCode:    404,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(content)),
		ContentLength: int64(len(content)),
	}
	return res
}
