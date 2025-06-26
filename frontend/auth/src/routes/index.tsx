import { MetaProvider, Title } from "@solidjs/meta";
import { useAction } from "@solidjs/router";
import { ModeToggle } from "~/components/mode-toggle";
import { authCheckAction } from "~/lib/auth";

export default function Home() {
  const authCheck = useAction(authCheckAction)
  authCheck()

  return (
    <>
      <MetaProvider>
        <Title>Home</Title>
      </MetaProvider>
      <main class="text-center mx-auto text-gray-700 p-4 bg-background">
        <p class="text-primary">
          Click here to change theme:  <ModeToggle />
        </p>
      </main>
    </>
  );
}
