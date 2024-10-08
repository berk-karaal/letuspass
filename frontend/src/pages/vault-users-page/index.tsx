import { retrieveVault } from "@/api/letuspass";
import { restQueryRetryFunc } from "@/common/queryRetry";
import { ActionIcon, Box, Group, Loader, Text, Title } from "@mantine/core";
import {
  IconArrowLeft,
  IconBriefcase2,
  IconUsersGroup,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router-dom";
import AddUserButtonAndModal from "./AddUserButtonAndModal";
import VaultUsersList from "./VaultUsersList";

function VaultUsersPage() {
  const { vaultId } = useParams();
  const navigate = useNavigate();

  const vaultQuery = useQuery({
    queryKey: ["vault", vaultId],
    queryFn: () => retrieveVault(Number(vaultId)),
    retry: restQueryRetryFunc,
  });

  return (
    <>
      <Group gap={"sm"} justify="flex-start">
        <ActionIcon
          variant="transparent"
          color="dark"
          onClick={() => navigate(-1)}
        >
          <IconArrowLeft size={"1.5rem"} />
        </ActionIcon>
        <IconBriefcase2 size={"1.75rem"} />
        <Title order={2} fw={"lighter"}>
          {vaultQuery.isSuccess ? (
            vaultQuery.data?.name
          ) : (
            <Loader color="gray" />
          )}
        </Title>
      </Group>

      <Group mt={"lg"} mb={"sm"} gap={"sm"} justify="flex-start">
        <IconUsersGroup size={"1.5rem"} />
        <Text fz={"h3"} fw={"lighter"}>
          Vault Users
        </Text>
        <Box ml={"auto"} mr={"xs"}>
          <AddUserButtonAndModal vaultId={Number(vaultId)} />
        </Box>
      </Group>

      <VaultUsersList vaultId={Number(vaultId)} />
    </>
  );
}

export default VaultUsersPage;
