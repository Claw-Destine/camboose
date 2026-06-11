import "bulma/css/bulma.min.css";
import { registerCambBody } from "./components/body/body";
import { registerNewProjectModal, registerProjectComponent, registerProjectsComponent } from "./components/projects/projects";
import { registerEditVersionModal, registerSpecsComponent, registerVersionItem } from "./components/specs/specs";

registerCambBody();
registerProjectComponent();
registerProjectsComponent();
registerNewProjectModal();
registerSpecsComponent();
registerVersionItem();
registerEditVersionModal();