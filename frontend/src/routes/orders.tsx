import {createFileRoute} from '@tanstack/react-router'
import {keepPreviousData, useQuery} from "@tanstack/react-query";
import {useState} from "react";

export const Route = createFileRoute('/orders')({
    component: Orders,
})

type Order = {
    _id: string;
    name: string;
};

function Orders() {
    const [page, setPage] = useState(1)

    const fetchOrders = (page = 1) =>
        fetch('/api/orders?page=' + page).then((response) => response.json())

    const { isPending, isError, data, error, isFetching, isPlaceholderData } = useQuery({
        queryKey: ['orders', page],
        queryFn: () => fetchOrders(page),
        placeholderData: keepPreviousData,
    })

    return (
        <div>
            {isPending ? (
                <div>Loading...</div>
            ) : isError ? (
                <div>Error: {error.message}</div>
            ) : (
                <div>
                    {data.orders.map((order: Order) => (
                        <p key={order._id}>{order.name}</p>
                    ))}
                </div>
            )}
            <span>Current page: {page}</span>
            <button
                onClick={() => setPage((old) => Math.max(old - 1, 1))}
                disabled={page === 1}
            >
                Previous Page
            </button>
            <button
                onClick={() => {
                    if (!isPlaceholderData && data.hasMore) {
                        setPage((old) => old + 1)
                    }
                }}
                // Disable the Next Page button until we know a next page is available
                disabled={isPlaceholderData || !data?.hasMore}
            >
                Next Page
            </button>
            {isFetching ? <span> Loading...</span> : null}
        </div>
    )
}
