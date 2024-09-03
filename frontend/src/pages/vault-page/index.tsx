import { retrieveVault } from "@/api/letuspass";
import { restQueryRetryFunc } from "@/common/queryRetry";
import { ActionIcon, Box, Group, Loader, Title } from "@mantine/core";
import {
  IconArrowLeft,
  IconBriefcase2,
  IconDotsVertical,
  IconFileTime,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router-dom";
import ThreeDotMenu from "./ThreeDotMenu";
import VaultItemList from "./VaultItemList";

function VaultPage() {
  const { vaultId } = useParams();
  const navigate = useNavigate();

  const vaultQuery = useQuery({
    queryKey: ["vault", Number(vaultId)],
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

        {vaultQuery.isSuccess && (
          <Box style={{ marginLeft: "auto" }}>
            <ActionIcon variant="transparent" color="dark" onClick={() => null}>
              <IconFileTime size={"1.5rem"} />
            </ActionIcon>
            <ThreeDotMenu
              vaultId={Number(vaultId)}
              vaultName={vaultQuery.data.name}
              target={
                <ActionIcon
                  variant="transparent"
                  color="dark"
                  onClick={() => null}
                >
                  <IconDotsVertical size={"1.5rem"} />
                </ActionIcon>
              }
            />
          </Box>
        )}
      </Group>
      <VaultItemList vaultId={Number(vaultId)} />
    </>
  );
}

export default VaultPage;
