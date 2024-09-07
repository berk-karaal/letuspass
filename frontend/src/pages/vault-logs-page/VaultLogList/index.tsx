import { listVaultAuditLogs } from "@/api/letuspass";
import { ControllersHandleVaultAuditLogsAuditLogResponseItem } from "@/api/letuspass.schemas";
import { Box, Group, Pagination, Stack, Text } from "@mantine/core";
import {
  IconAbc,
  IconBriefcase2,
  IconEdit,
  IconKey,
  IconLogout2,
  IconPlus,
  IconTrash,
  IconUserMinus,
  IconUserPlus,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import strftime from "strftime";
import classes from "./styles.module.css";

const ActionCategoryIconSize = "1.5rem";
const ActionIconSize = "1.5rem";

const ActionIconMappping = {
  vault_create: (
    <Group gap="0">
      <IconBriefcase2 size={ActionCategoryIconSize} />
      <IconPlus size={ActionIconSize} />
    </Group>
  ),
  vault_rename: (
    <Group gap="0">
      <IconBriefcase2 size={ActionCategoryIconSize} />
      <IconAbc size={ActionIconSize} />
    </Group>
  ),
  vault_delete: (
    <Group gap="0">
      <IconBriefcase2 size={ActionCategoryIconSize} />
      <IconTrash size={ActionIconSize} />
    </Group>
  ),
  vault_add_user: (
    <Group gap="0">
      <IconBriefcase2 size={ActionCategoryIconSize} />
      <IconUserPlus size={ActionIconSize} />
    </Group>
  ),
  vault_remove_user: (
    <Group gap="0">
      <IconBriefcase2 size={ActionCategoryIconSize} />
      <IconUserMinus size={ActionIconSize} />
    </Group>
  ),
  vault_user_left: (
    <Group gap="0">
      <IconBriefcase2 size={ActionCategoryIconSize} />
      <IconLogout2 size={ActionIconSize} />
    </Group>
  ),
  vault_item_create: (
    <Group gap="0">
      <IconKey size={ActionCategoryIconSize} />
      <IconPlus size={ActionIconSize} />
    </Group>
  ),
  vault_item_update: (
    <Group gap="0">
      <IconKey size={ActionCategoryIconSize} />
      <IconEdit size={ActionIconSize} />
    </Group>
  ),
  vault_item_delete: (
    <Group gap="0">
      <IconKey size={ActionCategoryIconSize} />
      <IconTrash size={ActionIconSize} />
    </Group>
  ),
};

function LogHighlight(text: string): React.ReactNode {
  return (
    <Text display={"inline"} c={"orange"}>
      {text}
    </Text>
  );
}

function VaultItemLink(
  title: string,
  vaultId: number,
  itemId: number
): React.ReactNode {
  return (
    <Text component={Link} to={`/app/vault/${vaultId}/item/${itemId}`}>
      {LogHighlight(`"${title}" (${itemId})`)}
    </Text>
  );
}

function auditLogToText(
  log: ControllersHandleVaultAuditLogsAuditLogResponseItem,
  vaultId: number
): React.ReactNode {
  switch (log.action_code) {
    case "vault_create":
      let create_data = log.action_data as { name: string };
      return (
        <Text>
          {LogHighlight(log.user.email)} created the vault with name{" "}
          {LogHighlight(`"${create_data.name}"`)}.
        </Text>
      );
    case "vault_rename":
      let rename_data = log.action_data as {
        old_name: string;
        new_name: string;
      };
      return (
        <Text>
          {LogHighlight(log.user.email)} renamed the vault{" "}
          {LogHighlight(`"${rename_data.old_name}"`)} to{" "}
          {LogHighlight(`"${rename_data.new_name}"`)}.
        </Text>
      );
    case "vault_delete":
      return <Text>{LogHighlight(log.user.email)} deleted the vault.</Text>;
    case "vault_add_user":
      let add_user_data = log.action_data as {
        added_user_email: string;
        permissions: string[];
      };
      return (
        <Text>
          {LogHighlight(log.user.email)} added user{" "}
          {LogHighlight(add_user_data.added_user_email)} with permissions:{" "}
          {LogHighlight(add_user_data.permissions.join(", "))}.
        </Text>
      );
    case "vault_remove_user":
      let remove_user_data = log.action_data as { removed_user_email: string };
      return (
        <Text>
          {LogHighlight(log.user.email)} removed user{" "}
          {LogHighlight(remove_user_data.removed_user_email)}.
        </Text>
      );
    case "vault_user_left":
      return <Text>{LogHighlight(log.user.email)} left vault.</Text>;
    case "vault_item_create":
      let create_item_data = log.action_data as { title: string };
      return (
        <Text>
          {LogHighlight(log.user.email)} creataed an item{" "}
          {VaultItemLink(
            create_item_data.title,
            vaultId,
            Number(log.vault_item?.id)
          )}
          .
        </Text>
      );
    case "vault_item_update":
      let update_item_data = log.action_data as { title: string };
      return (
        <Text>
          {LogHighlight(log.user.email)} edited an item{" "}
          {VaultItemLink(
            update_item_data.title,
            vaultId,
            Number(log.vault_item?.id)
          )}
          .
        </Text>
      );
    case "vault_item_delete":
      let delete_item_data = log.action_data as { title: string };
      return (
        <Text>
          {LogHighlight(log.user.email)} deleted an item{" "}
          {VaultItemLink(
            delete_item_data.title,
            vaultId,
            Number(log.vault_item?.id)
          )}
          .
        </Text>
      );
    default:
      return <Text>Unknown action: JSON.stringify(log);</Text>;
  }
}

function VaultLogBox({
  log,
  vaultId,
}: {
  log: ControllersHandleVaultAuditLogsAuditLogResponseItem;
  vaultId: number;
}) {
  return (
    <Box className={classes.VaultLogBox} py={"sm"} px={"md"}>
      <Group justify="space-between">
        {ActionIconMappping[log.action_code]}
        <Text size="sm">
          {strftime("%d-%m-%Y %H:%M %z", new Date(log.created_at))}
        </Text>
      </Group>
      <Text mt={"xs"}>{auditLogToText(log, vaultId)}</Text>
    </Box>
  );
}

export default function VaultLogList({ vaultId }: { vaultId: number }) {
  const PAGE_SIZE = 5;
  const [activePage, setActivePage] = useState(1);
  const [totalItemCount, setTotalItemCount] = useState(0);

  const [queryValues, setQueryValues] = useState({
    pageNumber: 1,
    pageSize: PAGE_SIZE,
  });

  const vaultLogsQuery = useQuery({
    queryKey: ["vaultLogs", vaultId, queryValues],
    queryFn: () =>
      listVaultAuditLogs(vaultId, {
        page: queryValues.pageNumber,
        page_size: queryValues.pageSize,
      }),
  });

  useEffect(() => {
    if (vaultLogsQuery.isSuccess && vaultLogsQuery.data.count >= 0) {
      setTotalItemCount(vaultLogsQuery.data.count);
    }
  }, [vaultLogsQuery.isFetching]);

  return (
    <>
      <Stack gap={"xs"}>
        {vaultLogsQuery.isSuccess &&
          vaultLogsQuery.data.results.map((log) => (
            <VaultLogBox key={log.id} log={log} vaultId={vaultId} />
          ))}
      </Stack>

      <Group mt={"lg"} justify={"center"}>
        <Pagination
          value={activePage}
          onChange={(value) => {
            setActivePage(value);
            setQueryValues((prev) => ({
              ...prev,
              pageNumber: value,
            }));
          }}
          total={Math.ceil(totalItemCount / PAGE_SIZE)}
        />
      </Group>
    </>
  );
}
