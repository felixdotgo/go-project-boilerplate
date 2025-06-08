import { A } from "@solidjs/router";
import { ModeToggle } from "~/components/mode-toggle";
import { Badge } from "~/components/ui/badge";
import { Button } from "~/components/ui/button";

export default function Home() {
  return (
    <main class="text-center mx-auto text-gray-700 p-4 bg-background">
      <Button variant="default">Button</Button>
      <Badge variant="error">Badge</Badge>
      <ModeToggle />
    </main>
  );
}
