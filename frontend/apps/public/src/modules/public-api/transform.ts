import i18next from "i18next";

import {
  canonicalProfile,
  canonicalProjects,
  canonicalResearchEntries,
} from "../profile-content";

import type {
  ProfileResponse,
  Project,
  ResearchEntry,
} from "./types";

type LocalizedText = {
  ja?: string | null;
  en?: string | null;
};

export type RawProfileResponse = {
  name?: LocalizedText;
  title?: LocalizedText;
  affiliation?: LocalizedText;
  lab?: LocalizedText;
  summary?: LocalizedText;
  skills?: LocalizedText[];
};

export type RawProject = {
  id: number;
  title?: LocalizedText;
  description?: LocalizedText;
  techStack?: string[];
  linkUrl?: string;
  year?: number;
};

export type RawResearchEntry = {
  id: number;
  title?: LocalizedText;
  summary?: LocalizedText;
  contentMd?: LocalizedText;
  year?: number;
};

function clone<T>(value: T): T {
  if (typeof structuredClone === "function") {
    return structuredClone(value);
  }
  return JSON.parse(JSON.stringify(value)) as T;
}

function selectLocalizedText(text?: LocalizedText | null): string | undefined {
  if (!text) {
    return undefined;
  }

  const language = i18next.language ?? "en";
  if (language.startsWith("ja") && text.ja) {
    return text.ja;
  }

  return text.en ?? text.ja ?? undefined;
}

export function transformProfile(
  raw: RawProfileResponse | undefined,
): ProfileResponse {
  if (raw && typeof (raw as unknown as ProfileResponse).name === "string") {
    return clone(raw as unknown as ProfileResponse);
  }

  const profile = clone(canonicalProfile);

  if (!raw) {
    return profile;
  }

  const name = selectLocalizedText(raw.name);
  if (name) {
    profile.name = name;
  }

  const headline = selectLocalizedText(raw.title);
  if (headline) {
    profile.headline = headline;
  }

  const summary = selectLocalizedText(raw.summary);
  if (summary) {
    profile.summary = summary;
  }

  const affiliation = selectLocalizedText(raw.affiliation);
  if (affiliation && profile.affiliations.length > 0) {
    profile.affiliations[0] = {
      ...profile.affiliations[0],
      organization: affiliation,
    };
  }

  const labName = selectLocalizedText(raw.lab);
  if (labName && profile.lab) {
    profile.lab = { ...profile.lab, name: labName };
  }

  if (raw.skills?.length) {
    const skillNames = raw.skills
      .map((skill) => selectLocalizedText(skill))
      .filter((name): name is string => Boolean(name));

    if (skillNames.length) {
      const engineeringGroup =
        profile.skillGroups.find(
          (group) => group.id === "software-engineering",
        ) ??
        profile.skillGroups[0] ??
        null;

      if (engineeringGroup) {
        const existingNames = new Set(
          engineeringGroup.items.map((item) => item.name),
        );
        engineeringGroup.items = [
          ...engineeringGroup.items,
          ...skillNames
            .filter((name) => !existingNames.has(name))
            .map((name, index) => ({
              id: `core-${index}`,
              name,
              level: "advanced" as const,
            })),
        ];
      }
    }
  }

  return profile;
}

const projectFallbackByTitle = new Map(
  canonicalProjects.map((project) => [project.title.toLowerCase(), project]),
);

const projectFallbackById = new Map(
  canonicalProjects.map((project, index) => [String(index + 1), project]),
);

export function transformProjects(
  projects: RawProject[] | undefined,
): Project[] {
  if (
    projects?.length &&
    typeof (projects[0] as unknown as Project).title === "string"
  ) {
    return clone(projects as unknown as Project[]);
  }

  if (!projects?.length) {
    return canonicalProjects.map(clone);
  }

  return projects.map((project) => {
    const localizedTitle = selectLocalizedText(project.title);
    const fallback =
      (localizedTitle
        ? projectFallbackByTitle.get(localizedTitle.toLowerCase())
        : undefined) ??
      projectFallbackById.get(String(project.id)) ??
      canonicalProjects[0];

    const result = clone(fallback);

    if (localizedTitle) {
      result.title = localizedTitle;
    }

    const description = selectLocalizedText(project.description);
    if (description) {
      result.description = description;
    }

    if (project.techStack?.length) {
      result.techStack = project.techStack;
    }

    if (project.linkUrl) {
      const existingRepo = result.links.findIndex(
        (link) => link.type === "repo",
      );
      if (existingRepo >= 0) {
        result.links[existingRepo] = {
          ...result.links[existingRepo],
          url: project.linkUrl,
        };
      } else {
        result.links = [
          ...result.links,
          { label: "Repository", url: project.linkUrl, type: "repo" },
        ];
      }
    }

    return result;
  });
}

const researchFallbackByTitle = new Map(
  canonicalResearchEntries.map((entry) => [entry.title.toLowerCase(), entry]),
);

const researchFallbackById = new Map(
  canonicalResearchEntries.map((entry, index) => [String(index + 1), entry]),
);

export function transformResearchEntries(
  entries: RawResearchEntry[] | undefined,
): ResearchEntry[] {
  if (
    entries?.length &&
    typeof (entries[0] as unknown as ResearchEntry).title === "string"
  ) {
    return clone(entries as unknown as ResearchEntry[]);
  }

  if (!entries?.length) {
    return canonicalResearchEntries.map(clone);
  }

  return entries.map((entry) => {
    const localizedTitle = selectLocalizedText(entry.title);
    const fallback =
      (localizedTitle
        ? researchFallbackByTitle.get(localizedTitle.toLowerCase())
        : undefined) ??
      researchFallbackById.get(String(entry.id)) ??
      canonicalResearchEntries[0];

    const result = clone(fallback);

    if (localizedTitle) {
      result.title = localizedTitle;
    }

    const summary = selectLocalizedText(entry.summary);
    if (summary) {
      result.summary = summary;
    }

    const contentMarkdown = selectLocalizedText(entry.contentMd);
    if (contentMarkdown) {
      result.contentMarkdown = contentMarkdown;
    }

    if (!result.contentHtml && contentMarkdown) {
      result.contentHtml = `<p>${contentMarkdown}</p>`;
    }

    return result;
  });
}
