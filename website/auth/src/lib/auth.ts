import Cookies from "js-cookie"
import { createServerCookie, getCookiesString, parseCookie } from "@solid-primitives/cookies"
import { getRequestEvent, isServer } from "solid-js/web";
import { action, redirect } from "@solidjs/router";

const authTokenKey = "sl_auth_token"

const [cookie, setCookie] = createServerCookie(authTokenKey)

const setAuthToken = (token: string) => {
  Cookies.set(authTokenKey, token)
};

const isAuthenticated = (): boolean => {
  let token : string | undefined = ""
  if (isServer) {
    const str = getCookiesString()
    token = parseCookie(str, authTokenKey)
  } else {
    token = Cookies.get(authTokenKey)
  }

  if (token && token !== undefined && token !== "" && token !== "undefined") {
    return true
  }

  return false
}

const authCheckAction = action(async () => {
  "use server"
  if (!isAuthenticated()) {
    throw redirect("/login")
  }
  return Promise.resolve()
})

export { setAuthToken, isAuthenticated, authCheckAction }
