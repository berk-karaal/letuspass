import { listVaultItems } from "@/api/letuspass";
import { Box, Group, LoadingOverlay, Pagination, Text } from "@mantine/core";
import { IconKey } from "@tabler/icons-react";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import classes from "./styles.module.css";

function VaultItemBox({
  itemTitle,
  vaultId,
  itemId,
}: {
  itemTitle: string;
  vaultId: number;
  itemId: number;
}) {
  return (
    <Box
      component={Link}
      to={`/app/vault/${vaultId}/item/${itemId}`}
      style={{
        textDecoration: "none",
        color: "inherit",
      }}
    >
      <Box className={classes.VaultItemBox} my={"xs"}>
        <Group px={"lg"} h={"100%"}>
          <IconKey size={"2rem"} />
          <Text size={"1.3rem"}>{itemTitle}</Text>
        </Group>
      </Box>
    </Box>
  );
}

export default function VaultItemList({ vaultId }: { vaultId: number }) {
  const PAGE_SIZE = 5;
  const [activePage, setPage] = useState(1);
  const [totalItemCount, setTotalItemCount] = useState(0);

  const vaultsQuery = useQuery({
    queryKey: ["vault", vaultId, "items", activePage],
    queryFn: () =>
      listVaultItems(vaultId, {
        page: activePage,
        page_size: PAGE_SIZE,
        ordering: "title",
      }),
    placeholderData: keepPreviousData,
    gcTime: 0,
  });

  useEffect(() => {
    if (vaultsQuery.isSuccess && vaultsQuery.data?.count > 0) {
      setTotalItemCount(vaultsQuery.data?.count);
    }
  }, [vaultsQuery]);

  return (
    <>
      {vaultsQuery.isSuccess && vaultsQuery.data?.count === 0 && (
        <Text ta={"center"} my={"lg"}>
          This vault is empty.
        </Text>
      )}
      <Group></Group>
      <Box pos={"relative"}>
        <LoadingOverlay
          visible={vaultsQuery.isFetching}
          zIndex={1000}
          overlayProps={{ radius: "sm", blur: 2 }}
        />
        {vaultsQuery.data?.results.map((vaultItem) => (
          <VaultItemBox
            key={vaultItem.id}
            itemId={Number(vaultItem.id)}
            itemTitle={vaultItem.title}
            vaultId={vaultId}
          />
        ))}
      </Box>
      <Group mt={"lg"} justify={"center"}>
        <Pagination
          value={activePage}
          onChange={setPage}
          total={Math.ceil(totalItemCount / PAGE_SIZE)}
        />
      </Group>
    </>
  );
}
