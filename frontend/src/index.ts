import "bulma/css/bulma.min.css";
import { registerCambBody } from "./components/body/body";
import { registerNewProjectModal, registerProjectComponent, registerProjectsComponent } from "./components/projects/projects";

registerCambBody();
registerProjectComponent();
registerProjectsComponent();
registerNewProjectModal();