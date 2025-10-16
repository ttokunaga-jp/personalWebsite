import { FormEvent, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { adminApi } from "./modules/admin-api";
import type {
  AdminProject,
  AdminResearch,
  AdminSummary,
  BlogPost,
  BlacklistEntry,
  Meeting,
  MeetingStatus
} from "./types";

const currentYear = new Date().getFullYear();

const emptyProjectForm = {
  titleJa: "",
  titleEn: "",
  descriptionJa: "",
  descriptionEn: "",
  techStack: "",
  linkUrl: "",
  year: String(currentYear),
  published: false,
  sortOrder: ""
};

const emptyResearchForm = {
  titleJa: "",
  titleEn: "",
  summaryJa: "",
  summaryEn: "",
  contentJa: "",
  contentEn: "",
  year: String(currentYear),
  published: false
};

const emptyBlogForm = {
  titleJa: "",
  titleEn: "",
  summaryJa: "",
  summaryEn: "",
  contentJa: "",
  contentEn: "",
  tags: "",
  published: false,
  publishedAt: ""
};

const emptyMeetingForm = {
  name: "",
  email: "",
  datetime: new Date().toISOString().slice(0, 16),
  durationMinutes: "30",
  meetUrl: "",
  status: "pending" as MeetingStatus,
  notes: ""
};

const emptyBlacklistForm = {
  email: "",
  reason: ""
};

function App() {
  const { t } = useTranslation();
  const [status, setStatus] = useState("unknown");
  const [summary, setSummary] = useState<AdminSummary | null>(null);
  const [projects, setProjects] = useState<AdminProject[]>([]);
  const [research, setResearch] = useState<AdminResearch[]>([]);
  const [blogs, setBlogs] = useState<BlogPost[]>([]);
  const [meetings, setMeetings] = useState<Meeting[]>([]);
  const [blacklist, setBlacklist] = useState<BlacklistEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [projectForm, setProjectForm] = useState({ ...emptyProjectForm });
  const [researchForm, setResearchForm] = useState({ ...emptyResearchForm });
  const [blogForm, setBlogForm] = useState({ ...emptyBlogForm });
  const [meetingForm, setMeetingForm] = useState({ ...emptyMeetingForm });
  const [blacklistForm, setBlacklistForm] = useState({ ...emptyBlacklistForm });

  const refreshAll = useCallback(async () => {
    setLoading(true);
    try {
      const [statusRes, summaryRes, projectRes, researchRes, blogRes, meetingRes, blacklistRes] =
        await Promise.all([
          adminApi.health(),
          adminApi.fetchSummary(),
          adminApi.listProjects(),
          adminApi.listResearch(),
          adminApi.listBlogs(),
          adminApi.listMeetings(),
          adminApi.listBlacklist()
        ]);

      setStatus(statusRes.data.status);
      setSummary(summaryRes.data);
      setProjects(projectRes.data);
      setResearch(researchRes.data);
      setBlogs(blogRes.data);
      setMeetings(meetingRes.data);
      setBlacklist(blacklistRes.data);
      setError(null);
    } catch (err) {
      console.error(err);
      setError("status.error");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refreshAll();
  }, [refreshAll]);

  const run = useCallback(
    async (operation: () => Promise<unknown>) => {
      try {
        await operation();
        await refreshAll();
      } catch (err) {
        console.error(err);
        setError("status.error");
      }
    },
    [refreshAll]
  );

  const handleCreateProject = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = {
      title: { ja: projectForm.titleJa, en: projectForm.titleEn },
      description: { ja: projectForm.descriptionJa, en: projectForm.descriptionEn },
      techStack: projectForm.techStack.split(",").map((item) => item.trim()).filter(Boolean),
      linkUrl: projectForm.linkUrl.trim(),
      year: Number(projectForm.year) || currentYear,
      published: projectForm.published,
      sortOrder: projectForm.sortOrder === "" ? null : Number(projectForm.sortOrder)
    };

    await run(async () => {
      await adminApi.createProject(payload);
      setProjectForm({ ...emptyProjectForm });
    });
  };

  const handleCreateResearch = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = {
      title: { ja: researchForm.titleJa, en: researchForm.titleEn },
      summary: { ja: researchForm.summaryJa, en: researchForm.summaryEn },
      contentMd: { ja: researchForm.contentJa, en: researchForm.contentEn },
      year: Number(researchForm.year) || currentYear,
      published: researchForm.published
    };

    await run(async () => {
      await adminApi.createResearch(payload);
      setResearchForm({ ...emptyResearchForm });
    });
  };

  const handleCreateBlog = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const publishedAt = blogForm.publishedAt ? new Date(blogForm.publishedAt).toISOString() : null;
    const payload = {
      title: { ja: blogForm.titleJa, en: blogForm.titleEn },
      summary: { ja: blogForm.summaryJa, en: blogForm.summaryEn },
      contentMd: { ja: blogForm.contentJa, en: blogForm.contentEn },
      tags: blogForm.tags.split(",").map((tag) => tag.trim()).filter(Boolean),
      published: blogForm.published,
      publishedAt
    };

    await run(async () => {
      await adminApi.createBlog(payload);
      setBlogForm({ ...emptyBlogForm });
    });
  };

  const handleCreateMeeting = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = {
      name: meetingForm.name,
      email: meetingForm.email,
      datetime: new Date(meetingForm.datetime).toISOString(),
      durationMinutes: Number(meetingForm.durationMinutes) || 30,
      meetUrl: meetingForm.meetUrl,
      status: meetingForm.status,
      notes: meetingForm.notes
    };

    await run(async () => {
      await adminApi.createMeeting(payload);
      setMeetingForm({ ...emptyMeetingForm, datetime: new Date().toISOString().slice(0, 16) });
    });
  };

  const handleCreateBlacklist = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = {
      email: blacklistForm.email.trim(),
      reason: blacklistForm.reason.trim()
    };

    await run(async () => {
      await adminApi.createBlacklist(payload);
      setBlacklistForm({ ...emptyBlacklistForm });
    });
  };

  const toggleProjectPublished = (project: AdminProject) =>
    run(async () => {
      await adminApi.updateProject(project.id, {
        title: project.title,
        description: project.description,
        techStack: project.techStack,
        linkUrl: project.linkUrl,
        year: project.year,
        published: !project.published,
        sortOrder: project.sortOrder ?? null
      });
    });

  const toggleResearchPublished = (item: AdminResearch) =>
    run(async () => {
      await adminApi.updateResearch(item.id, {
        title: item.title,
        summary: item.summary,
        contentMd: item.contentMd,
        year: item.year,
        published: !item.published
      });
    });

  const toggleBlogPublished = (post: BlogPost) =>
    run(async () => {
      await adminApi.updateBlog(post.id, {
        title: post.title,
        summary: post.summary,
        contentMd: post.contentMd,
        tags: post.tags,
        published: !post.published,
        publishedAt: post.publishedAt ?? null
      });
    });

  const updateMeetingStatus = (meeting: Meeting, status: MeetingStatus) =>
    run(async () => {
      await adminApi.updateMeeting(meeting.id, {
        name: meeting.name,
        email: meeting.email,
        datetime: meeting.datetime,
        durationMinutes: meeting.durationMinutes,
        meetUrl: meeting.meetUrl,
        status,
        notes: meeting.notes
      });
    });

  const deleteProject = (id: number) => run(() => adminApi.deleteProject(id));
  const deleteResearch = (id: number) => run(() => adminApi.deleteResearch(id));
  const deleteBlog = (id: number) => run(() => adminApi.deleteBlog(id));
  const deleteMeeting = (id: number) => run(() => adminApi.deleteMeeting(id));
  const deleteBlacklistEntry = (id: number) => run(() => adminApi.deleteBlacklist(id));

  return (
    <div className="min-h-screen bg-slate-100">
      <header className="bg-slate-900 p-6 text-white">
        <h1 className="text-2xl font-bold">{t("dashboard.title")}</h1>
        <p className="text-sm text-slate-300">{t("dashboard.subtitle")}</p>
      </header>
      <main className="mx-auto flex max-w-6xl flex-col gap-6 p-6">
        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-800">{t("dashboard.systemStatus")}</h2>
          <div className="mt-4 flex flex-wrap gap-6">
            <div className="flex-1 min-w-[220px] rounded-md bg-slate-900 p-4 text-white">
              <span className="font-mono uppercase tracking-wide text-slate-400">
                {t("dashboard.apiStatus")}
              </span>
              <p className="text-2xl font-bold text-emerald-400">{status}</p>
            </div>
            {summary && (
              <div className="flex-1 min-w-[220px] rounded-md border border-slate-200 bg-white p-4">
                <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500">
                  {t("summary.title")}
                </h3>
                <ul className="mt-2 space-y-1 text-sm text-slate-700">
                  <li>Projects: {summary.publishedProjects} published / {summary.draftProjects} draft</li>
                  <li>Research: {summary.publishedResearch} published / {summary.draftResearch} draft</li>
                  <li>Blogs: {summary.publishedBlogs} published / {summary.draftBlogs} draft</li>
                  <li>Pending meetings: {summary.pendingMeetings}</li>
                  <li>Blacklist entries: {summary.blacklistEntries}</li>
                </ul>
              </div>
            )}
          </div>
        </section>

        {loading && (
          <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
            <p className="text-slate-600">{t("status.loading")}</p>
          </section>
        )}

        {!loading && error && (
          <section className="rounded-lg border border-rose-200 bg-white p-6 shadow-sm">
            <p className="text-rose-600">{t(error)}</p>
          </section>
        )}

        {!loading && !error && (
          <>
            <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
              <header className="mb-4 flex items-center justify-between">
                <h2 className="text-lg font-semibold text-slate-800">{t("projects.title")}</h2>
              </header>
              <form className="grid gap-3 md:grid-cols-2" onSubmit={handleCreateProject}>
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Title (ja)"
                  value={projectForm.titleJa}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, titleJa: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Title (en)"
                  value={projectForm.titleEn}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, titleEn: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Description (ja)"
                  value={projectForm.descriptionJa}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, descriptionJa: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Description (en)"
                  value={projectForm.descriptionEn}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, descriptionEn: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Tech stack (comma separated)"
                  value={projectForm.techStack}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, techStack: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Link URL"
                  value={projectForm.linkUrl}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, linkUrl: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Year"
                  value={projectForm.year}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, year: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Sort order"
                  value={projectForm.sortOrder}
                  onChange={(event) => setProjectForm((prev) => ({ ...prev, sortOrder: event.target.value }))}
                />
                <label className="flex items-center gap-2 text-sm text-slate-700">
                  <input
                    checked={projectForm.published}
                    onChange={(event) => setProjectForm((prev) => ({ ...prev, published: event.target.checked }))}
                    type="checkbox"
                  />
                  Published
                </label>
                <button className="rounded bg-slate-900 px-4 py-2 text-white" type="submit">
                  {t("actions.create")}
                </button>
              </form>

              <div className="mt-6 space-y-4">
                {projects.map((project) => (
                  <div key={project.id} className="rounded border border-slate-200 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-semibold text-slate-800">{project.title.ja || project.title.en}</h3>
                        <p className="text-sm text-slate-500">{project.year}</p>
                      </div>
                      <div className="flex gap-2">
                        <button
                          className="rounded border border-slate-300 px-3 py-1 text-sm"
                          type="button"
                          onClick={() => toggleProjectPublished(project)}
                        >
                          {project.published ? "Unpublish" : "Publish"}
                        </button>
                        <button
                          className="rounded border border-rose-200 bg-rose-50 px-3 py-1 text-sm text-rose-600"
                          type="button"
                          onClick={() => deleteProject(project.id)}
                        >
                          {t("actions.delete")}
                        </button>
                      </div>
                    </div>
                    <p className="mt-2 text-sm text-slate-600">{project.description.ja || project.description.en}</p>
                    <p className="mt-1 text-xs text-slate-500">Stack: {project.techStack.join(", ")}</p>
                  </div>
                ))}
              </div>
            </section>

            <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
              <header className="mb-4 flex items-center justify-between">
                <h2 className="text-lg font-semibold text-slate-800">{t("research.title")}</h2>
              </header>
              <form className="grid gap-3 md:grid-cols-2" onSubmit={handleCreateResearch}>
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Title (ja)"
                  value={researchForm.titleJa}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, titleJa: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Title (en)"
                  value={researchForm.titleEn}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, titleEn: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Summary (ja)"
                  value={researchForm.summaryJa}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, summaryJa: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Summary (en)"
                  value={researchForm.summaryEn}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, summaryEn: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Content (ja)"
                  value={researchForm.contentJa}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, contentJa: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Content (en)"
                  value={researchForm.contentEn}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, contentEn: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Year"
                  value={researchForm.year}
                  onChange={(event) => setResearchForm((prev) => ({ ...prev, year: event.target.value }))}
                />
                <label className="flex items-center gap-2 text-sm text-slate-700">
                  <input
                    checked={researchForm.published}
                    onChange={(event) => setResearchForm((prev) => ({ ...prev, published: event.target.checked }))}
                    type="checkbox"
                  />
                  Published
                </label>
                <button className="rounded bg-slate-900 px-4 py-2 text-white" type="submit">
                  {t("actions.create")}
                </button>
              </form>
              <div className="mt-6 space-y-4">
                {research.map((item) => (
                  <div key={item.id} className="rounded border border-slate-200 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-semibold text-slate-800">{item.title.ja || item.title.en}</h3>
                        <p className="text-sm text-slate-500">{item.year}</p>
                      </div>
                      <div className="flex gap-2">
                        <button
                          className="rounded border border-slate-300 px-3 py-1 text-sm"
                          type="button"
                          onClick={() => toggleResearchPublished(item)}
                        >
                          {item.published ? "Unpublish" : "Publish"}
                        </button>
                        <button
                          className="rounded border border-rose-200 bg-rose-50 px-3 py-1 text-sm text-rose-600"
                          type="button"
                          onClick={() => deleteResearch(item.id)}
                        >
                          {t("actions.delete")}
                        </button>
                      </div>
                    </div>
                    <p className="mt-2 text-sm text-slate-600">{item.summary.ja || item.summary.en}</p>
                  </div>
                ))}
              </div>
            </section>

            <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
              <header className="mb-4 flex items-center justify-between">
                <h2 className="text-lg font-semibold text-slate-800">{t("blogs.title")}</h2>
              </header>
              <form className="grid gap-3 md:grid-cols-2" onSubmit={handleCreateBlog}>
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Title (ja)"
                  value={blogForm.titleJa}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, titleJa: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Title (en)"
                  value={blogForm.titleEn}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, titleEn: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Summary (ja)"
                  value={blogForm.summaryJa}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, summaryJa: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Summary (en)"
                  value={blogForm.summaryEn}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, summaryEn: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Content (ja)"
                  value={blogForm.contentJa}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, contentJa: event.target.value }))}
                />
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Content (en)"
                  value={blogForm.contentEn}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, contentEn: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Tags (comma separated)"
                  value={blogForm.tags}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, tags: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  type="datetime-local"
                  value={blogForm.publishedAt}
                  onChange={(event) => setBlogForm((prev) => ({ ...prev, publishedAt: event.target.value }))}
                />
                <label className="flex items-center gap-2 text-sm text-slate-700">
                  <input
                    checked={blogForm.published}
                    onChange={(event) => setBlogForm((prev) => ({ ...prev, published: event.target.checked }))}
                    type="checkbox"
                  />
                  Published
                </label>
                <button className="rounded bg-slate-900 px-4 py-2 text-white" type="submit">
                  {t("actions.create")}
                </button>
              </form>

              <div className="mt-6 space-y-4">
                {blogs.map((post) => (
                  <div key={post.id} className="rounded border border-slate-200 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-semibold text-slate-800">{post.title.ja || post.title.en}</h3>
                        {post.publishedAt && (
                          <p className="text-xs text-slate-500">Published at {new Date(post.publishedAt).toLocaleString()}</p>
                        )}
                      </div>
                      <div className="flex gap-2">
                        <button
                          className="rounded border border-slate-300 px-3 py-1 text-sm"
                          type="button"
                          onClick={() => toggleBlogPublished(post)}
                        >
                          {post.published ? "Unpublish" : "Publish"}
                        </button>
                        <button
                          className="rounded border border-rose-200 bg-rose-50 px-3 py-1 text-sm text-rose-600"
                          type="button"
                          onClick={() => deleteBlog(post.id)}
                        >
                          {t("actions.delete")}
                        </button>
                      </div>
                    </div>
                    <p className="mt-2 text-sm text-slate-600">{post.summary.ja || post.summary.en}</p>
                    <p className="mt-1 text-xs text-slate-500">Tags: {post.tags.join(", ")}</p>
                  </div>
                ))}
              </div>
            </section>

            <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
              <header className="mb-4 flex items-center justify-between">
                <h2 className="text-lg font-semibold text-slate-800">{t("meetings.title")}</h2>
              </header>
              <form className="grid gap-3 md:grid-cols-2" onSubmit={handleCreateMeeting}>
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Name"
                  value={meetingForm.name}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, name: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Email"
                  value={meetingForm.email}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, email: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  type="datetime-local"
                  value={meetingForm.datetime}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, datetime: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Duration (minutes)"
                  value={meetingForm.durationMinutes}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, durationMinutes: event.target.value }))}
                />
                <input
                  className="rounded border border-slate-300 p-2"
                  placeholder="Meet URL"
                  value={meetingForm.meetUrl}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, meetUrl: event.target.value }))}
                />
                <select
                  className="rounded border border-slate-300 p-2"
                  value={meetingForm.status}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, status: event.target.value as MeetingStatus }))}
                >
                  <option value="pending">Pending</option>
                  <option value="confirmed">Confirmed</option>
                  <option value="cancelled">Cancelled</option>
                </select>
                <textarea
                  className="rounded border border-slate-300 p-2 md:col-span-2"
                  placeholder="Notes"
                  value={meetingForm.notes}
                  onChange={(event) => setMeetingForm((prev) => ({ ...prev, notes: event.target.value }))}
                />
                <button className="rounded bg-slate-900 px-4 py-2 text-white" type="submit">
                  {t("actions.create")}
                </button>
              </form>

              <div className="mt-6 space-y-4">
                {meetings.map((meeting) => (
                  <div key={meeting.id} className="rounded border border-slate-200 p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-semibold text-slate-800">{meeting.name}</h3>
                        <p className="text-sm text-slate-500">{meeting.email}</p>
                        <p className="text-xs text-slate-500">{new Date(meeting.datetime).toLocaleString()}</p>
                      </div>
                      <div className="flex gap-2">
                        <select
                          className="rounded border border-slate-300 p-1 text-sm"
                          value={meeting.status}
                          onChange={(event) => updateMeetingStatus(meeting, event.target.value as MeetingStatus)}
                        >
                          <option value="pending">Pending</option>
                          <option value="confirmed">Confirmed</option>
                          <option value="cancelled">Cancelled</option>
                        </select>
                        <button
                          className="rounded border border-rose-200 bg-rose-50 px-3 py-1 text-sm text-rose-600"
                          type="button"
                          onClick={() => deleteMeeting(meeting.id)}
                        >
                          {t("actions.delete")}
                        </button>
                      </div>
                    </div>
                    {meeting.notes && <p className="mt-2 text-sm text-slate-600">{meeting.notes}</p>}
                  </div>
                ))}
              </div>
            </section>

            <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
              <header className="mb-4 flex items-center justify-between">
                <h2 className="text-lg font-semibold text-slate-800">{t("blacklist.title")}</h2>
              </header>
              <form className="flex flex-col gap-3 md:flex-row" onSubmit={handleCreateBlacklist}>
                <input
                  className="flex-1 rounded border border-slate-300 p-2"
                  placeholder="Email"
                  value={blacklistForm.email}
                  onChange={(event) => setBlacklistForm((prev) => ({ ...prev, email: event.target.value }))}
                />
                <input
                  className="flex-1 rounded border border-slate-300 p-2"
                  placeholder="Reason"
                  value={blacklistForm.reason}
                  onChange={(event) => setBlacklistForm((prev) => ({ ...prev, reason: event.target.value }))}
                />
                <button className="rounded bg-slate-900 px-4 py-2 text-white" type="submit">
                  {t("actions.create")}
                </button>
              </form>

              <div className="mt-6 space-y-3">
                {blacklist.map((entry) => (
                  <div key={entry.id} className="flex items-center justify-between rounded border border-slate-200 p-3">
                    <div>
                      <p className="font-medium text-slate-800">{entry.email}</p>
                      {entry.reason && <p className="text-sm text-slate-500">{entry.reason}</p>}
                    </div>
                    <button
                      className="rounded border border-rose-200 bg-rose-50 px-3 py-1 text-sm text-rose-600"
                      type="button"
                      onClick={() => deleteBlacklistEntry(entry.id)}
                    >
                      {t("actions.delete")}
                    </button>
                  </div>
                ))}
              </div>
            </section>
          </>
        )}
      </main>
    </div>
  );
}

export default App;
