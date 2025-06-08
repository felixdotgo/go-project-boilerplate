import { MetaProvider, Title } from "@solidjs/meta";
import { ModeToggle } from "~/components/mode-toggle";

export default function Home() {
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
