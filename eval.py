min_tut_diff = 25
min_tut_list = []
rec_mapping = {}
for i in range(1000):
    with open(f'data/{i+1}/stats.txt') as f:
        lines = f.readlines()
        tut_values = [int(line[-3:-1]) for line in lines[6:22]]
        rec_values = [int(line[-3:-1]) for line in lines[24:39]]
        rec_mapping[i+1] = max(rec_values) - min(rec_values)
        diff = max(tut_values) - min(tut_values)
        if diff == min_tut_diff:
            min_tut_list.append(i+1)
        elif diff < min_tut_diff:
            min_tut_diff = diff
            min_tut_list = [i+1]

min_rec_diff = 25 
best_schedule = 0
for schedule in min_tut_list:
    if rec_mapping[schedule] < min_rec_diff:
        min_rec_diff = rec_mapping[schedule]
        best_schedule = schedule

print("Minimum tutorial difference:",min_tut_diff, "\nMinimum rec difference:", min_rec_diff, "\nSchedule:", best_schedule)