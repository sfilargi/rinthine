#!/usr/bin/python3 -u

from github import Github

def load_token():
        with open('/home/admin/.github_token') as f:
                return f.readline().strip('\n\r')

def main():
        g = Github(load_token())
        r = g.get_repo('rinthine-com/rinthine')
        pulls = []
        prs = r.get_pulls()
        for pr in prs:
            if not pr.mergeable:
                continue
            thumb_ups = 0
            for r in pr.as_issue().get_reactions():
                if r.content == '+1':
                    thumb_ups +1
                elif r.content == '-1':
                    thumb_ups -= 1
            pulls.append((pr, thumb_ups))
        if len(pulls) == 0:
            return
        pull = sorted(pulls, reverse = True, key = lambda x: x[1])[0][0]
        pull.merge(merge_method='squash')

if __name__ == '__main__':
        main()
