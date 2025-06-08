import { A } from "@solidjs/router";
import { MetaProvider, Title } from "@solidjs/meta";
import { Button } from "~/components/ui/button";

export default function NotFound() {
  return (
    <>
      <MetaProvider>
        <Title>Not Found</Title>
      </MetaProvider>
      <div class="mx-auto p-5 w-full text-center">
        <h1 class="font-bold text-4xl">404</h1>
        <p>Page not found.</p>
        <Button as={A} variant="link" href="/" class="mt-4">Click here to go back to Home</Button>
      </div>
    </>
  );
}
