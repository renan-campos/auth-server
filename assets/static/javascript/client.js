export class BackendClient {
  constructor() {
    this.root = window.location.href;
  }

    async SendOtp(otp, handler) {
        Post(this.root + "token", otp, handler)
    }

}

// helpers {
function Post(route, body, handler) {
    fetch(route, 
        {method: "POST", headers: { "Content-Type": "application/text" }, body: body}
    ).then((resp) => { handler(resp) })
}
// } helpers
