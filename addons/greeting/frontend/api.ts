import { api, ProblemError, throwProblem } from "@margince/frontend/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

// Every call this extension's screens make, so each page reads the same keys.
const notesKey = ["ext", "acme-greeting", "notes"] as const;

// One deal's notes, or the most recent across deals.
export function useAcmeNotes(dealId?: string) {
  return useQuery({
    queryKey: [...notesKey, dealId ?? "recent"],
    queryFn: async () => {
      const { data, error } = await api.GET("/ext/acme-greeting/notes", {
        params: { query: dealId ? { deal_id: dealId } : {} },
      });
      if (error) throwProblem(error);
      return data.notes;
    },
  });
}

export function useAddAcmeNote(dealId: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: async (body: string) => {
      const { data, error } = await api.POST("/ext/acme-greeting/notes", { body: { deal_id: dealId, body } });
      if (error) throwProblem(error);
      return data;
    },
    onSuccess: () => client.invalidateQueries({ queryKey: ["ext", "acme-greeting"] }),
  });
}

// Core's GET /deals/{id}/score, answered by the scorer this extension provides.
export function useDealScore(dealId: string) {
  return useQuery({
    queryKey: ["ext", "acme-greeting", "score", dealId],
    queryFn: async () => {
      const { data, error } = await api.GET("/deals/{id}/score", { params: { path: { id: dealId } } });
      if (error) throwProblem(error);
      return data;
    },
  });
}

// Whether a call was refused because the reader's role does not hold this
// extension's RBAC object, which an admin grants per role.
export function refusedForRole(error: unknown): boolean {
  if (!(error instanceof ProblemError)) return false;
  const problem = error.problem;
  return typeof problem === "object" && problem !== null && "status" in problem && problem.status === 403;
}
