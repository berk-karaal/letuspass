import { listVaults } from "@/api/letuspass";
import { Box, Group, LoadingOverlay, Pagination, Text } from "@mantine/core";
import { IconBriefcase2 } from "@tabler/icons-react";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import classes from "./styles.module.css";

function VaultBox({
  vaultName,
  vaultId,
}: {
  vaultName: string;
  vaultId: number;
}) {
  return (
    <Box
      component={Link}
      to={`/app/vault/${vaultId}`}
      style={{
        textDecoration: "none",
        color: "inherit",
      }}
    >
      <Box className={classes.VaultBox} my={"xs"}>
        <Group px={"lg"} h={"100%"}>
          <IconBriefcase2 size={"2rem"} />
          <Text size={"1.3rem"}>{vaultName}</Text>
        </Group>
      </Box>
    </Box>
  );
}

function VaultList() {
  const PAGE_SIZE = 5;
  const [activePage, setPage] = useState(1);
  const [totalItemCount, setTotalItemCount] = useState(0);

  const vaultsQuery = useQuery({
    queryKey: ["vaults", activePage],
    queryFn: () =>
      listVaults({ page: activePage, page_size: PAGE_SIZE, ordering: "name" }),
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
          You don't have any vaults yet.
        </Text>
      )}
      <Box pos={"relative"}>
        <LoadingOverlay
          visible={vaultsQuery.isFetching}
          zIndex={1000}
          overlayProps={{ radius: "sm", blur: 2 }}
        />
        {vaultsQuery.data?.results.map((vault) => (
          <VaultBox key={vault.id} vaultName={vault.name} vaultId={vault.id} />
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

export default VaultList;
