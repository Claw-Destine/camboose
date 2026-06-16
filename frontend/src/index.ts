import "bulma/css/bulma.min.css";
import { registerNewProjectModal, registerProjectComponent, registerProjectsComponent } from "./components/projects/projects";
import { registerEditVersionModal, registerSpecsComponent, registerVersionItem } from "./components/specs/specs";
import "./components/body/body"

registerProjectComponent();
registerProjectsComponent();
registerNewProjectModal();
registerSpecsComponent();
registerVersionItem();
registerEditVersionModal();