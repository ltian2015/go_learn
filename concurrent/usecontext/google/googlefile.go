// Package google provides a function to do Google searches using the Google Web
// Search API. See https://developers.google.com/web-search/docs/
// This package is an example to accompany https://blog.golang.org/context.
// It is not intended for use by others.
// Google has since disabled its search API,
// and so this package is no longer useful.
// google 包提供了一个使用 Google Web 搜索 API 进行谷歌搜索的函数。
// 参见 https://developers.google.com/web-search/docs/
// 本包作为 https://blog.golang.org/context 的配套示例。
// 不建议其他人使用。 Google 已经禁用了其搜索 API，
// 因此本包已不再有用。
/** Client 是一个 HTTP 客户端。其零值（[DefaultClient]）是一个使用 了[DefaultTransport]的可用客户端，。
	[Client.Transport] 通常具有内部状态（缓存的 TCP 连接），因此应重用 Client，而不是按需创建。
	 Client 可以安全地被多个 goroutine 并发使用。
	 Client 比 [RoundTripper]（如 [Transport]）更高级，还处理如 cookies 和重定向等 HTTP 细节。
	 在伴随（following）重定向时，Client 会转发（forward）在初始请求[Request] 上设置的所有headers信息，
	除了以下两种情况：
   - 当将敏感的headers信息如 "Authorization"、"WWW-Authenticate" 和 "Cookie" 转发到不受信任目标时。
     当重定向到的域不是初始域的子域或完全匹配时，这些headers会被忽略。
     例如，从 "foo.com" 重定向到 "foo.com" 或 "sub.foo.com" 时会转发敏感headers，
     但重定向到 "bar.com" 时不会转发敏感的headers。
	- 当转发带有非 nil 的 cookie Jar 的 "Cookie"  header时。
	因为每次重定向都可能改变 cookie jar 的状态，重定向可能会更改初始请求中设置的 cookie。
	转发 "Cookie" header时，任何被更改的 cookie 都会被省略，除非Jar 会插入这些已更新值的 cookie（假设源匹配）。
	如果 Jar 为 nil，则初始 cookie 会原样转发。
    -----------
	Do() 发送 HTTP 请求并返回 HTTP 响应，遵循客户端配置的策略（如重定向、cookie、认证）。
	如果是由客户端策略（如 CheckRedirect）或 HTTP 对话失败（如网络连接问题）引起，则会返回错误。
	非 2xx 状态代码不会导致错误（2xx状态码表示错误）。
	 如果返回的错误为零，[Response] 将包含一个非零的 Body，用户应关闭该 Body。
	 如果 Body 没有被读取到 EOF 并关闭，[客户端]的底层[RoundTripper]（通常是[Transport]）可能
	 无法为后续的 "keep-alive "请求重新使用与服务器的持久 TCP 连接。
	请求正文（如果非nil）将由底层传输系统关闭，即使在出错时也是如此。
	请求正文可在 Do() 返回后异步关闭。
	---------------
	结构体http.Request有个名为ctx的私有属性,可通过其Context()方法访问.
	WithContext 返回req 的一个浅拷贝，并将其上下文属性更改为 ctx。
	但要求提供的 ctx 必须是非 nil 的。
	 对于传出的客户端请求，上下文控制请求及其响应的整个生命周期：获取连接、发送请求以及读取响应头和正文。
	 要创建一个带有上下文的新请求，请使用 [NewRequestWithContext]。
	 要对带有新上下文的请求进行深拷贝，请使用 [Request.Clone]。
    **/

package google

import (
	"context"
	"encoding/json"
	"net/http"

	"com.example/golearn/concurrent/usecontext/userip"
)

// Results is an ordered list of search results.
// Reusults 是一个有序的搜索结果列表。
type Results []Result

// A Result contains the title and URL of a search result.
// Result 包含搜索结果的标题和 URL。
type Result struct {
	Title, URL string
}

// Search sends query to Google search and returns the results.
// Search 向 Google 搜索发送查询并返回结果。
func Search(ctx context.Context, query string) (Results, error) {
	// Prepare the Google Search API request.
	// 准备 Google 搜索 API 请求。
	req, err := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)

	// 如果 ctx 携带了用户 IP 地址，则将其转发到服务器。
	// Google APIs 使用用户 IP 来区分服务器发起的请求与最终用户请求。
	// 如果 ctx 携带了用户 IP 地址，则将其转发到服务器。
	if userIP, ok := userip.FromContext(ctx); ok {
		q.Set("userip", userIP.String())
	}
	req.URL.RawQuery = q.Encode()

	// Issue the HTTP request and handle the response. The httpDo function
	// cancels the request if ctx.Done is closed.
	// 发出 HTTP 请求并处理响应。httpDo 函数在 ctx.Done 关闭时取消请求。
	var results Results
	responseHandler := func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Parse the JSON search result.
		// https://developers.google.com/web-search/docs/#fonje
		// 解析 JSON 搜索结果。
		var data struct {
			ResponseData struct {
				Results []struct {
					TitleNoFormatting string
					URL               string
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			println(err)
			return err
		}
		for _, res := range data.ResponseData.Results {
			results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
		}
		return nil
	}
	err = httpDo(ctx, req, responseHandler)
	// httpDo waits for the closure we provided to return, so it's safe to
	// read results here.
	return results, err
}

// httpDo 发出 HTTP 请求并调用 f 处理响应。
// 如果 ctx.Done 在请求或 f 运行时关闭，httpDo 将取消请求，等待 f 退出，并返回 ctx.Err。
// 否则，httpDo 返回 f 的错误。
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {

	chErr := make(chan error, 1)

	req = req.WithContext(ctx)

	go func() {
		chErr <- f(http.DefaultClient.Do(req))
	}() // 在 goroutine 中运行 HTTP 请求，并将响应传递给 f进行处理，f是响应处理函数。

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chErr:
		return err
	}
}
