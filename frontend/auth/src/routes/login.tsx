import { MetaProvider, Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import { BsApple, BsFacebook, BsGoogle } from "solid-icons/bs";
import { ModeToggle } from "~/components/mode-toggle";
import { Button } from "~/components/ui/button";
import { TextField, TextFieldInput, TextFieldLabel } from "~/components/ui/text-field";
import { showToast } from "~/components/ui/toast";
import { setAuthToken } from "~/lib/auth";

export default function Login() {
  const navigator = useNavigate();

  function handleLogin() {
    showToast({ title: "Login successful", variant: "success" });
    setAuthToken("test-token");
    navigator("/", { replace: true });
  }

  return (
    <>
      <MetaProvider>
        <Title>Login</Title>
      </MetaProvider>
      <div class="flex w-full mx-auto">
        <div class="hidden lg:block w-1/2 h-screen bg-primary"></div>
        <div class="relative w-full h-screen min-h-96 lg:w-1/2 justify-items-center content-center">
          <div class="fixed top-4 right-4">
            <ModeToggle />
          </div>
          <div class="w-full max-w-sm p-4">
            <TextField class="mb-4">
              <TextFieldLabel>Email</TextFieldLabel>
              <TextFieldInput type="text" />
            </TextField>

            <TextField class="mb-4">
              <TextFieldLabel>Password</TextFieldLabel>
              <TextFieldInput type="password" />
            </TextField>

            <Button class="w-full mb-4" onClick={handleLogin}>Login</Button>

            <div class="relative mb-4">
              <div class="absolute inset-0 flex items-center">
                <span class="w-full border-t" />
              </div>
              <div class="relative flex justify-center text-xs uppercase">
                <span class="bg-background px-2 text-muted-foreground">Or continue with</span>
              </div>
            </div>

            <div class="flex flex-auto justify-between">
              <Button class="w-1/3 mr-1" variant={"outline"}><BsGoogle /> Google</Button>
              <Button class="w-1/3 mr-1" variant={"outline"}><BsFacebook /> Facebook</Button>
              <Button class="w-1/3" variant={"outline"}><BsApple /> Apple</Button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
