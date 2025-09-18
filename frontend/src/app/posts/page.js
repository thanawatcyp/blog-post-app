import PostsList from "./PostsList";

export default function PostsPage() {
  // Provide empty initial data, PostsList will fetch data on client side
  const initialData = {
    items: [],
    total: 0,
    page: 1,
    page_size: 10
  };

  return <PostsList initialData={initialData} />;
}
