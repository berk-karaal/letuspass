import { Group, Title } from "@mantine/core";
import CreateVaultButtonAndModal from "./CreateVaultButtonAndModal";
import VaultList from "./VaultList";

function AppPage() {
  return (
    <>
      <Group mb={"sm"} justify={"space-between"}>
        <Title fw={"lighter"}>Vaults</Title>
        <CreateVaultButtonAndModal />
      </Group>
      <VaultList />
    </>
  );
}

export default AppPage;
