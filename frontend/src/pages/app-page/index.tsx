import { Button, Group, Title } from "@mantine/core";
import { IconPlus } from "@tabler/icons-react";
import VaultList from "./VaultList";

function AppPage() {
  return (
    <>
      <Group mb={"sm"} justify={"space-between"}>
        <Title fw={"lighter"}>Vaults</Title>
        <Button leftSection={<IconPlus size={"20px"} />}>New Vault</Button>
      </Group>
      <VaultList />
    </>
  );
}

export default AppPage;
