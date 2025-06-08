import { MetaProvider, Title } from "@solidjs/meta";
import { ModeToggle } from "~/components/mode-toggle";
import { Button } from "~/components/ui/button";
import { Separator } from "~/components/ui/separator";
import { TextField, TextFieldInput, TextFieldLabel } from "~/components/ui/text-field";

export default function Login() {
  return (
    <>
      <MetaProvider>
        <Title>Login</Title>
      </MetaProvider>
      <div class="flex w-full mx-auto">
        <div class="hidden lg:block w-1/2 h-screen bg-primary"></div>
        <div class="relative w-full lg:w-1/2 justify-items-center content-center">
          <div class="fixed top-4 right-4">
            <ModeToggle />
          </div>
          <div class="w-full max-w-sm">
            <TextField class="mb-4">
              <TextFieldLabel>Email</TextFieldLabel>
              <TextFieldInput type="text" />
            </TextField>

            <TextField class="mb-4">
              <TextFieldLabel>Password</TextFieldLabel>
              <TextFieldInput type="password" />
            </TextField>

            <Button class="w-full mb-4">Login</Button>

            <div class="relative mb-4">
              <div class="absolute inset-0 flex items-center">
                <span class="w-full border-t" />
              </div>
              <div class="relative flex justify-center text-xs uppercase">
                <span class="bg-background px-2 text-muted-foreground">Or continue with</span>
              </div>
            </div>

            <Button class="w-full" variant={"secondary"}>Google</Button>
          </div>
        </div>
      </div>
    </>
  )
}
