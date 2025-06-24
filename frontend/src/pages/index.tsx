import Sidebar from "@/components/Sidebar/Sidebar";
import AppView from "@/components/views/AppView";
import { Route } from "@/types/GlobalTypes";

export default function Home({
  route,
  setRoute,
}: {
  route: Route;
  setRoute: (route: Route) => void;
}) {
  return <>
    <Sidebar activeRoute={route} setRoute={setRoute} />
    <AppView route={route} setRoute={setRoute} />
  </>
}
