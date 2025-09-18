"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { authenticatedFetch, logout } from "@/lib/auth";

export default function PostsList({ initialData }) {
  const [posts, setPosts] = useState(initialData.items || []);
  const [page, setPage] = useState(initialData.page || 1);
  const [totalPages, setTotalPages] = useState(
    Math.ceil(initialData.total / initialData.page_size) || 1
  );
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [createFormData, setCreateFormData] = useState({
    title: "",
    content: "",
    author: ""
  });
  const [createLoading, setCreateLoading] = useState(false);
  const router = useRouter();

  const fetchPosts = async (pageNum = 1, query = "") => {
    try {
      setLoading(true);
      setError("");

      const params = new URLSearchParams({
        page: pageNum.toString(),
        page_size: "10",
      });
      if (query.trim()) params.append("q", query.trim());

      console.log('page number change : ' + pageNum);

      const response = await authenticatedFetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/posts?${params}`
      );

      if (!response.ok) throw new Error("Failed to fetch posts");

      const data = await response.json();
      setPosts(data.items || []);
      setTotalPages(Math.ceil(data.total / data.page_size));
      setPage(data.page);
    } catch (err) {
      if (err.message !== "Authentication required") {
        setError(err.message || "Failed to load posts");
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // Fetch new data when page changes (skip initial load)
    if (page !== (initialData.page || 1)) {
      fetchPosts(page, searchQuery);
    }
  }, [page]);

  useEffect(() => {
    // Fetch data on component mount if we have empty initial data
    if (initialData.items.length === 0) {
      fetchPosts(1, '');
    }
  }, []);

  const handleSearch = (e) => {
    e.preventDefault();
    fetchPosts(1, searchQuery);
  };

  const handleLogout = async () => {
    await logout();
  };

  const handleCreatePost = async (e) => {
    e.preventDefault();

    if (!createFormData.title.trim() || !createFormData.content.trim() || !createFormData.author.trim()) {
      setError("Please fill in all fields");
      return;
    }

    try {
      setCreateLoading(true);
      setError("");

      const response = await authenticatedFetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/posts/create`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(createFormData),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to create post");
      }

      // Show success alert
      alert("🎉 Post created successfully!");

      // Reset form and close
      setCreateFormData({ title: "", content: "", author: "" });
      setShowCreateForm(false);

      // Refresh posts list
      fetchPosts(1, searchQuery);
    } catch (err) {
      // Show error alert popup
      alert("❌ " + (err.message || "Failed to create post"));
      setError(err.message || "Failed to create post");
    } finally {
      setCreateLoading(false);
    }
  };

  const formatDate = (dateString) =>
    new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center py-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Blog Posts</h1>
            <p className="mt-1 text-sm text-gray-600">
              Discover and read amazing content
            </p>
          </div>
          <div className="flex gap-3">
            <button
              onClick={() => setShowCreateForm(true)}
              className="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition"
            >
              Create Post
            </button>
            <button
              onClick={handleLogout}
              className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md transition"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Create Post Modal */}
        {showCreateForm && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
              <div className="p-6">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="text-2xl font-bold text-gray-900">Create New Post</h2>
                  <button
                    onClick={() => {
                      setShowCreateForm(false);
                      setCreateFormData({ title: "", content: "", author: "" });
                      setError("");
                    }}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>

                <form onSubmit={handleCreatePost} className="space-y-4">
                  <div>
                    <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                      Title
                    </label>
                    <input
                      type="text"
                      id="title"
                      value={createFormData.title}
                      onChange={(e) => setCreateFormData({ ...createFormData, title: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                      placeholder="Enter post title..."
                      required
                    />
                  </div>

                  <div>
                    <label htmlFor="author" className="block text-sm font-medium text-gray-700 mb-1">
                      Author
                    </label>
                    <input
                      type="text"
                      id="author"
                      value={createFormData.author}
                      onChange={(e) => setCreateFormData({ ...createFormData, author: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                      placeholder="Enter author name..."
                      required
                    />
                  </div>

                  <div>
                    <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-1">
                      Content
                    </label>
                    <textarea
                      id="content"
                      value={createFormData.content}
                      onChange={(e) => setCreateFormData({ ...createFormData, content: e.target.value })}
                      rows="6"
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                      placeholder="Write your post content..."
                      required
                    />
                  </div>

                  <div className="flex gap-3 pt-4">
                    <button
                      type="submit"
                      disabled={createLoading}
                      className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition"
                    >
                      {createLoading ? "Creating..." : "Create Post"}
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        setShowCreateForm(false);
                        setCreateFormData({ title: "", content: "", author: "" });
                        setError("");
                      }}
                      className="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-700 font-medium py-2 px-4 rounded-md transition"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        )}

        {/* Search */}
        <form onSubmit={handleSearch} className="flex gap-4 mb-6">
          <input
            type="text"
            placeholder="Search posts..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
          />
          <button
            type="submit"
            disabled={loading}
            className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-md transition"
          >
            Search
          </button>
        </form>

        {/* Error */}
        {error && (
          <div className="mb-6 bg-red-50 border border-red-400 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}

        {/* Posts */}
        {loading && posts.length === 0 ? (
          <p className="text-center text-gray-500">Loading posts...</p>
        ) : posts.length === 0 ? (
          <p className="text-center text-gray-500">No posts found</p>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {posts.map((post) => (
              <article
                key={post.id}
                className="bg-white rounded-lg shadow p-6 hover:shadow-md transition"
              >
                <h2 className="text-xl font-semibold mb-2">{post.title}</h2>
                <p className="text-gray-600 mb-3 line-clamp-3">{post.content}</p>
                <div className="text-sm text-gray-500 flex justify-between">
                  <span>By {post.author || "Anonymous"}</span>
                  <time>{formatDate(post.created_at)}</time>
                </div>
              </article>
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-8 flex justify-center items-center gap-4">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1 || loading}
              className="px-4 py-2 border rounded disabled:opacity-50"
            >
              Previous
            </button>
            <span>
              Page {page} of {totalPages}
            </span>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages || loading}
              className="px-4 py-2 border rounded disabled:opacity-50"
            >
              Next
            </button>
          </div>
        )}
      </main>
    </div>
  );
}
