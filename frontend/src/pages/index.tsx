import AppView from "@/components/views/AppView";
import { Route } from "@/types/GlobalTypes";

export default function Home({
  route,
  setRoute,
}: {
  route: Route;
  setRoute: (route: Route) => void;
}) {
  return <AppView route={route} setRoute={setRoute} />;
}
