import { AddNewLicensePageView } from "./AddNewLicensePageView";

export default {
	title: "pages/DeploySettingsPage/AddNewLicensePageView",
	component: AddNewLicensePageView,
};

export const Default = {
	args: {
		isSavingLicense: false,
		didSavingFailed: false,
	},
};
